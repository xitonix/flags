package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.EnableAutoEnv()
	str := flags.DefaultBucket.StringP("name", "n", "default", "usage of name").WithEnv("Custom")
	_ = flags.DefaultBucket.StringP("something-longer", "", "default", "usage of name is a bit longer").Var()
	flags.Parse()
	fmt.Println(*str)
}
