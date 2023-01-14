package conf

import (
	"fmt"
	"testing"
)

func TestLoadFromJsonFile(t *testing.T) {
	root := LoadFromJsonFile("./test.json")
	fmt.Println(root)
}
