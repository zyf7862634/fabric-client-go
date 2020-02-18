package utils

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	var snowflake *Snowflake
	snowflake, err := NewSnowflake(1)
	if err != nil {
		t.Fatal(err)
	}

	var temp []string
	for i := 0; i < 10; i++ {
		id := snowflake.Generate()
		temp = append(temp, id.String())
	}
	ret := RemoveDuplicate(temp)
	fmt.Printf("total %02d, first: %s\n", len(ret), ret[0])
}
