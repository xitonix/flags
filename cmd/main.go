package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.String("name", "usage of name").WithDefault("abc")
	_ = flags.StringP("something-longer", `usage of name is a bit longer`, "s").MarkAsDeprecated().Var()
	flags.Parse()
	fmt.Println(*str)
}
