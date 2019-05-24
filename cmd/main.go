package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetEnvPrefix("ALG")
	flags.EnableAutoEnv()
	str := flags.String("name", "usage of name").WithShort("n").WithEnv("A_NAME")
	_ = flags.StringD("something-longer", `usage of name is a bit longer`, "Johnny").Var()
	flags.Parse()
	fmt.Println(*str)
}
