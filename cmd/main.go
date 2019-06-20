package main

import (
	"fmt"
	"net"
	"time"

	"go.xitonix.io/flags"
)

func main() {

	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	str := flags.StringP("flag", "usage of name", "A").WithKey("ABC").WithDefault("abc").Var()
	i := flags.Int("int-flag", "usage of int flag").WithKey("IntK")
	_ = flags.StringP("x-flag", `usage of name is a bit longer`, "x").Var()
	dur := flags.DurationP("ttl", "duration usage", "t").WithDefault(5 * time.Second)
	ss := flags.StringSliceP("colours", "Colour pallet", "c").WithDefault([]string{"Pink"})
	ip := flags.IPAddressP("ip-address", "IP address usage", "i").WithValidRange(net.ParseIP("10.10.10.10"))
	flags.Parse()
	fmt.Println("Value:", *str)
	fmt.Println("Int:", i.Get())
	fmt.Printf("Dur:%v\n", dur.Get())
	fmt.Printf("Colours:%v\n", ss.Get())
	fmt.Printf("IP:%v\n", ip.Get())
}
