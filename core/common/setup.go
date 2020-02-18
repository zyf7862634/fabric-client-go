package common

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"sort"
	"sync"
)

var logger = SetupModuleLogger("core.common")

var (
	setupOnce sync.Once
	setup     *FabricSetup
)

type FabricSetup struct {
	//配置文件只有一个
	ConfigFile       string
	OrgAdmin         string
	OrgName          string
	OrgID            string
	ChannelID        string
	OrdererEndPoint  string
	Affiliation      string
	CertCaName       string
	IdentityTypeUser string

	Sdk           *fabsdk.FabricSDK //实例化后的sdk
	MspClient     *msp.Client
	resMgmtClient *resmgmt.Client //资源客户端
	Client        *channel.Client //通道客户端
}

func GetFabricSetupIns() *FabricSetup {
	setupOnce.Do(func() {
		svrCfg := GetSvrConfigIns()
		setup = &FabricSetup{
			ConfigFile:       svrCfg.GetCfgFile(FabricCliConfig),
			OrgAdmin:         FabricAdminUser,
			OrgName:          svrCfg.GetCfgString(FabricOrgName),
			ChannelID:        svrCfg.GetCfgString(FabricOrgChannel),
			OrdererEndPoint:  svrCfg.GetCfgString(FabricOrderer),
			Affiliation:      svrCfg.GetCfgString(FabricAffiliation),
			CertCaName:       svrCfg.GetCfgString(FabricCertCaName),
			IdentityTypeUser: FabricIdentityTypeUser,
		}
		setup.setupSDK()
		setup.setupMspClient()
		setup.setupResMgmtClient()
		setup.setupChannelClient()
	})
	return setup
}

func (fs *FabricSetup) GetUserNameFromCert(cert string) (string, error) {
	if !GetSvrConfigIns().GetCfgBool(ServerUserCertOrName) {
		return cert, nil
	}

	decoded, _ := pem.Decode([]byte(cert))
	if decoded == nil {
		return "", fmt.Errorf("failed to decode cert: [%s]", cert)
	}

	certificate, err := x509.ParseCertificate(decoded.Bytes)
	if err != nil || certificate == nil || certificate.Subject.CommonName == "" {
		return "", fmt.Errorf("failed to parse certificate: %v, %v", err, certificate)
	}

	return certificate.Subject.CommonName, nil
}

func (fs *FabricSetup) RegisteredAndEnrollUser(userName string) (string, bool) {
	if userName == fs.OrgAdmin {
		return fs.signingIdentity(userName)
	}

	// Register the new user
	enrollmentSecret, err := setup.MspClient.Register(
		&msp.RegistrationRequest{
			Name:        userName,
			Type:        fs.IdentityTypeUser,
			Affiliation: fs.Affiliation,
		})
	// err if user is already enrolled
	if err == nil {
		// Enroll the new user
		err = fs.MspClient.Enroll(userName, msp.WithSecret(enrollmentSecret))
		if err != nil {
			return err.Error(), false
		}
	}

	if err := fs.joinChannel(userName); err != nil {
		logger.Errorf("join channel failed: %v", err)
	}

	return fs.signingIdentity(userName)
}

func (fs *FabricSetup) GetPeerList() ([]string, error) {
	configBackend, _ := fs.Sdk.Config()
	endpointConfig, _ := fab.ConfigFromBackend(configBackend)
	networkPeers := endpointConfig.NetworkPeers()

	peers := make([]string, 0, len(networkPeers))
	for _, p := range networkPeers {
		peers = append(peers, p.URL)
	}

	sort.Strings(peers)
	return peers, nil
}

func (fs *FabricSetup) signingIdentity(userName string) (string, bool) {
	// Get the new user's signing identity
	si, err := fs.MspClient.GetSigningIdentity(userName)
	if err != nil {
		return err.Error(), false
	}

	certInfo := si.EnrollmentCertificate()
	return string(certInfo), true
}

func (fs *FabricSetup) joinChannel(userName string) error {
	orgUserClientContext := fs.Sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(fs.OrgName))
	resMgmtClient, err := resmgmt.New(orgUserClientContext)
	if err != nil {
		return errors.WithMessage(err, "Failed to create new resource management channelClient")
	}

	options := []resmgmt.RequestOption{
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	}
	if err = resMgmtClient.JoinChannel(setup.ChannelID, options...); err != nil {
		return err
	}

	return nil
}

func (fs *FabricSetup) setupSDK() {
	sdk, err := fabsdk.New(config.FromFile(fs.ConfigFile))
	if err != nil {
		logger.Fatalf("failed to create new SDK: %v", err)
	}
	fs.Sdk = sdk
	logger.Infof("Fabric sdk created")
}

func (fs *FabricSetup) setupResMgmtClient() {
	options := []fabsdk.ContextOption{fabsdk.WithUser(fs.OrgAdmin), fabsdk.WithOrg(fs.OrgName)}
	resourceManagerClientContext := fs.Sdk.Context(options...)
	client, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		logger.Fatalf("fail to create channel management client from Admin identity: %v", err)
	}
	setup.resMgmtClient = client
	logger.Infof("Resource management client created")
}

func (fs *FabricSetup) setupMspClient() {
	ctxProvider := fs.Sdk.Context()
	fs.MspClient, _ = msp.New(ctxProvider, msp.WithOrg(fs.OrgName))

	/*registrarEnrollID, registrarEnrollSecret := fs.getRegistrarEnrollmentCredentials(ctxProvider)
	err := fs.MspClient.Enroll(registrarEnrollID, msp.WithSecret(registrarEnrollSecret))
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}*/
}

func (fs *FabricSetup) getRegistrarEnrollmentCredentials(ctxProvider context.ClientProvider) (string, string) {
	ctx, err := ctxProvider()
	if err != nil {
		logger.Fatalf("failed to get context: %v", err)
	}

	clientConfig := ctx.IdentityConfig().Client()
	caOrg := clientConfig.Organization
	logger.Debugf("CaOrg is %s", caOrg)
	caConfig, exist := ctx.IdentityConfig().CAConfig(caOrg)
	if !exist {
		logger.Fatalf("CAConfig failed: %v\n", err)
	}

	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}

func (fs *FabricSetup) setupChannelClient() {
	options := []fabsdk.ContextOption{fabsdk.WithUser(fs.OrgAdmin), fabsdk.WithOrg(fs.OrgName)}
	clientChannelContext := fs.Sdk.ChannelContext(fs.ChannelID, options...)

	var err error
	fs.Client, err = channel.New(clientChannelContext)
	if err != nil {
		logger.Fatalf("create channel channelClient error: %v", err)
	}
	logger.Infof("Channel channelClient created")
}
