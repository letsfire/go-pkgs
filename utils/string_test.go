package utils

import (
	"fmt"
	"testing"
)

func TestSSliceDiff(t *testing.T) {
	ss1 := []string{"1", "2", "3", "4"}
	ss2 := []string{"1", "2", "3", "4", "5"}
	fmt.Println(SSliceDiff(ss1, ss2))
	fmt.Println(SSliceDiff(ss2, ss1))
}
