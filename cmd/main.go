package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.String("a-flag", "usage of name").WithKey("ABC").WithDefault("abc")
	_ = flags.String("x-flag", `usage of name is a bit longer`).Var()
	flags.Parse()
	fmt.Println(*str)
}
