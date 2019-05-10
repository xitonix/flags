package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetEnvPrefix("AAA")
	str := flags.StringPD("name", "n", "focus", "usage of name").WithEnv("ENV_NAME")
	_ = flags.StringD("something-longer", "hh", `usage of name is a bit longer`).Var()
	flags.Parse()
	fmt.Println(*str)
}
