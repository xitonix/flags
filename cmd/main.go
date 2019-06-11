package main

import (
	"fmt"

	"go.xitonix.io/flags"
)

func main() {

	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.StringP("flag", "usage of name", "A").WithKey("ABC").WithDefault("abc").Var()
	i := flags.Int("int-flag", "usage of int flag").WithDefault(8000).WithKey("IntK")
	_ = flags.StringP("x-flag", `usage of name is a bit longer`, "x").Var()
	flags.Parse()
	fmt.Println("Value:", *str)
	fmt.Println("Int:", i.Get())
}
