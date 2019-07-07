package main

import (
	"fmt"

	"github.com/xitonix/flags"
	"github.com/xitonix/flags/by"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	flags.SetSortOrder(by.RequiredFirst)
	i := flags.Int("int-flag", "usage of int flag").WithShort("i")
	j := flags.Int("j-flag", "usage of int flag").WithShort("j")
	b := flags.Bool("bool", `usage of name is a bit longer`).WithShort("b")
	f := flags.Float64("float", `usage of name is a bit longer`).WithShort("f")
	is := flags.IntSlice("is", "int slice flag").WithShort("s")

	flags.Parse()
	fmt.Println("Bool:", b.Get())
	fmt.Println("I:", i.Get())
	fmt.Println("J:", j.Get())
	fmt.Println("F:", f.Get())
	fmt.Println("Slice:", is.Get())
}
