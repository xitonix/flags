package main

import (
	"fmt"
	"go.xitonix.io/flags"
)

func main() {
	flags.SetEnvPrefix("AAA")
	str := flags.String("NAME", "usage of name").WithEnv("ENV_NAME").WithShort("--c")
	_ = flags.StringD("something-longer", `usage of name is a bit longer`, "default value").Var()
	flags.Parse()
	fmt.Println(*str)
}
