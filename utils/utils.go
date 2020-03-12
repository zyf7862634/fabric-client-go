package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const RegexpJsonReplace = "(\"\\w+\":|\"|\\{|\\})"

func ElapsedTime(fun string, s time.Time) string {
	elapsed := time.Since(s)
	return fmt.Sprintf("func %s elapsed: %v", fun, elapsed)
}

func GetCurrentExeFileDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetArrayArgBytes(jsonArgs string) [][]byte {
	data := make([][]byte, 0, 1)
	array := strings.Split(jsonArgs, ",")
	for _, v := range array {
		data = append(data, []byte(v))
	}
	return data
}

func GetJsonArgBytes(jsonArgs string) [][]byte {
	data := make([][]byte, 0, 1)
	data = append(data, []byte(jsonArgs))
	return data
}

func NotSupportFunctionError(funcName string) (interface{}, error) {
	return "Failed", fmt.Errorf("not support function name '%s'", funcName)
}

func SuccessFunctionMessage() (interface{}, error) {
	return "Success", nil
}
