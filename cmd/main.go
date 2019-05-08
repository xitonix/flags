package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	str := flags.DefaultBucket.StringP("name", "n", "default", "usage of name")
	flags.Parse()
	fmt.Println(str.Env().Name(), str.Get())

}
