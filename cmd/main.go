package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.String("name", "usage of name").WithKey("ABC").WithDefault("abc")
	_ = flags.String("something-longer", `usage of name is a bit longer`).WithKey("abc").Var()
	flags.Parse()
	fmt.Println(*str)
}
