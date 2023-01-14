package utils

import "fmt"

func Println(args ...interface{}) {
	fmt.Println("----------------------------------------------------------------------------------")
	for i := range args {
		fmt.Printf("%#v", args[i])
		fmt.Println("\n----------------------------------------------------------------------------------")
	}
}
