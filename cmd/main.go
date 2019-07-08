package main

import (
	"fmt"

	"github.com/xitonix/flags"
)

func main() {
	port := flags.Int("port-number", "Port number").WithShort("p")
	log := flags.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()
	flags.Parse()
	fmt.Println("Port", port.Get())
	fmt.Println("Log", *log)
}
