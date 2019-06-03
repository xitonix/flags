package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {

	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.StringP("flag", "usage of name", "A").WithKey("ABC").WithDefault("abc").Var()
	_ = flags.StringP("x-flag", `usage of name is a bit longer`, "x").Var()
	flags.Parse()
	fmt.Println("Value:", *str)
}
