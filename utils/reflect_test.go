package utils

import (
	"fmt"
	"testing"
)

type User struct {
	Name   string `test:"name,1,1,2"`
	Gender int    `test:"gender"`
	year   int    `test:"year"`
}

func TestFieldsNameByTag(t *testing.T) {
	fmt.Println(FieldsNameByTag(User{}, "test", true))
}

func BenchmarkFieldsNameByTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FieldsNameByTag(User{}, "test", false)
	}
}
