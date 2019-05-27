package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()

	str := flags.StringP("a-flag", "usage of name", "B").WithKey("ABC").WithDefault("abc")
	_ = flags.StringP("x-flag", `usage of name is a bit longer`, "A").Var()
	flags.Parse()
	fmt.Println(*str)
}
