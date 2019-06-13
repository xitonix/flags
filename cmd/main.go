package main

import (
	"fmt"
	"time"

	"go.xitonix.io/flags"
)

func main() {

	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.StringP("flag", "usage of name", "A").WithKey("ABC").WithDefault("abc").Var()
	i := flags.Int("int-flag", "usage of int flag").WithDefault(8000).WithKey("IntK")
	_ = flags.StringP("x-flag", `usage of name is a bit longer`, "x").Var()
	dur := flags.DurationP("ttl", "duration usage", "t").WithDefault(5 * time.Second)
	ss := flags.StringSliceP("colours", "Colour pallet", "c").WithDefault([]string{"Pink"})
	flags.Parse()
	fmt.Println("Value:", *str)
	fmt.Println("Int:", i.Get())
	fmt.Printf("Dur:%v\n", dur.Get())
	fmt.Printf("Colours:%v\n", ss.Get())
}
