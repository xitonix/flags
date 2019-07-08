package main

import (
	"errors"
	"fmt"

	"github.com/xitonix/flags"
	"github.com/xitonix/flags/by"
)

func main() {
	// Enabling auto key generation
	flags.EnableAutoKeyGeneration()
	flags.SetKeyPrefix("PFX")

	// Customising deprecation and required marks
	flags.SetDeprecationMark("**DEPRECATED**")
	flags.SetRequiredFlagMark("**")

	// Changing the sort order
	flags.SetSortOrder(by.LongNameDescending)

	// With long and short forms
	port := flags.Int("port-number", "Port number").WithShort("p")

	// With default value
	log := flags.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()

	// Input Validation
	weekend := flags.StringSlice("weekends", "Weekends").WithValidRange(true, "Sat, Sun").WithTrimming()
	numRange := flags.Int8("number", "A flag with validation callback").
		WithValidationCallback(func(in int8) error {
			if in > 10 {
				return errors.New("--number must be less than 10")
			}
			return nil
		})

	// CIDR and IP address
	net := flags.CIDR("network", "Network definition. Example 192.168.1.1/16")
	endpoint := flags.IPAddress("endpoint", "The IP address of the remote server")

	// Deprecated flag
	_ = flags.Int("port", "Legacy port number. Use -p, --port-number instead").MarkAsDeprecated()

	// Required flag
	rate := flags.Float64("rate", "Conversion rate").Required()

	// Hidden flag
	hidden := flags.Bool("enable-logging", "Secret flag").Hide()

	t := flags.Time("start-time", "Start time")
	ttl := flags.Duration("ttl", "Time to live")

	// Counter flags
	counter := flags.Counter("counter", "Repeat counter")
	verbose := flags.Verbosity("Verbosity. Repeat -v for higher verbosity levels. Example -vv")

	flags.Parse()
	fmt.Println("Port", port.Get())
	fmt.Println("Log", *log)
	fmt.Println("Weekend", weekend.Get())
	fmt.Println("Network", net.Get())
	fmt.Println("Endpoint", endpoint.Get())
	fmt.Println("Rate", rate.Get())
	fmt.Println("Hidden", hidden.Get())
	fmt.Println("Range", numRange.Get())
	fmt.Println("Time", t.Get())
	fmt.Println("TTL", ttl.Get())
	fmt.Println("Counter", counter.Get())
	fmt.Println("Verbosity", verbose.Get())

	for _, flag := range flags.DefaultBucket.Flags() {
		fmt.Printf("--%s (%s) %s\n", flag.LongName(), flag.Type(), flag.Usage())
	}
}
