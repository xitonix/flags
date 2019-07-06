package main

import (
	"fmt"
	"net"
	"time"

	"github.com/xitonix/flags"
	"github.com/xitonix/flags/by"
)

func main() {
	flags.SetKeyPrefix("ALG")
	flags.EnableAutoKeyGeneration()
	flags.SetSortOrder(by.RequiredFirst)
	str := flags.String("flag", "usage of name").WithKey("ABC").WithDefault("abc").Var()
	i := flags.Int("int-flag", "usage of int flag").WithKey("IntK")
	_ = flags.String("x-flag", `usage of name is a bit longer`).Var()
	dur := flags.Duration("ttl", "duration usage").WithDefault(5 * time.Second)
	ss := flags.StringSlice("colours", "Colour pallet").WithDefault([]string{"Pink"})
	ip := flags.IPAddress("ip-address", "IP address usage").WithValidRange(net.ParseIP("10.10.10.10"))
	sl := flags.StringSlice("ss", "string slice").WithValidRange(false, "A", "B")
	cid := flags.CIDR("network", "network usage").Required()
	mp := flags.StringMap("mappings", "String map flag").Required()
	flags.Parse()
	fmt.Println("Value:", *str)
	fmt.Println("Int:", i.Get())
	fmt.Printf("Dur:%v\n", dur.Get())
	fmt.Printf("Colours:%v\n", ss.Get())
	fmt.Printf("IP:%v\n", ip.Get())
	fmt.Printf("SL:%v\n", sl.Get())
	fmt.Printf("CID:%v\n", cid.Get().String())
	fmt.Printf("MAP:%v\n", mp.Get())
}
