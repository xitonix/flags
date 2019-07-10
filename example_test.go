package flags_test

import (
	"errors"
	"fmt"

	"github.com/xitonix/flags"
	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/config"
	"github.com/xitonix/flags/core"
)

func Example() {
	port := flags.Int("port-number", "Port number").WithShort("p")
	log := flags.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()

	flags.Parse()

	fmt.Println("Port", port.Get())
	fmt.Println("Log", *log)
}

func Example_bucket() {
	bucket := flags.NewBucket(config.WithAutoKeys(), config.WithSortOrder(by.KeyAscending))
	port := bucket.Int("port-number", "Port number").WithShort("p")
	log := bucket.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()

	bucket.Parse()

	fmt.Println("Port", port.Get())
	fmt.Println("Log", *log)
}

func Example_bucketCallbacks() {

	preCallback := func(flag core.Flag, value string) error {
		fmt.Printf("%s will be set to %s\n", flag.LongName(), value)
		return nil
	}

	postCallback := func(flag core.Flag, value string) error {
		fmt.Printf("%s has been set to %s\n", flag.LongName(), value)
		return nil
	}

	bucket := flags.NewBucket(config.WithPreSetCallback(preCallback), config.WithPostSetCallback(postCallback))
	external := bucket.Int("internal", "Internal port number").WithShort("e")
	internal := bucket.Int("external", "External port number").WithShort("i")

	bucket.Parse()

	fmt.Println("External Port", external.Get())
	fmt.Println("Internal Port", internal.Get())
}

func Example_deprecated() {
	port := flags.Int("port-number", "Port number").WithDefault(8080).WithShort("p")
	_ = flags.Int("port", "Legacy port number. Use -p, --port-number instead").MarkAsDeprecated()

	flags.Parse()
	fmt.Println("Port", port.Get())
}

func Example_validationCallback() {
	port := flags.Int("port-number", "Port number").WithValidationCallback(func(in int) error {
		if in > 9000 {
			return errors.New("The port number must be less than 9000")
		}
		return nil
	})

	flags.Parse()
	fmt.Println("Port", port.Get())
}

func Example_validationRange() {
	port := flags.Int("port-number", "Port number").WithDefault(8080).WithValidRange(8080, 8081)
	flags.Parse()
	fmt.Println("Port", port.Get())
}

func Example_slice() {
	numbers := flags.IntSlice("numbers", "Numbers list").WithShort("N")
	flags.Parse()
	fmt.Println("Numbers", numbers.Get())
	// ./binary -N "1, 2, 3"
}

func Example_withKeys() {
	flags.EnableAutoKeyGeneration()
	flags.SetKeyPrefix("PFX")
	numbers := flags.IntSlice("numbers", "Numbers list")
	log := flags.String("log-file", "The path to the log file").WithKey("LOG")
	flags.Parse()
	fmt.Println("Numbers", numbers.Get())
	fmt.Println("Log", log.Get())
	// PFX_PORT_NUMBER=80 PFX_LOG=/var/log/output.log ./binary
}

func Example_customSource() {
	flags.EnableAutoKeyGeneration()

	ms := flags.NewMemorySource()
	ms.Add("PORT_NUMBER", "9090")
	flags.PrependSource(ms)

	port := flags.Int("port-number", "Port number").WithShort("p")
	flags.Parse()
	fmt.Println("Port", port.Get())
}

func Example_hiddenFlag() {
	port := flags.Int("port-number", "Port number").WithShort("p")
	// Hidden flags won't be printed in the help output
	rate := flags.Float32("rate", "Secret conversion rate").Hide()

	flags.Parse()

	fmt.Println("Port", port.Get())
	fmt.Println("Rate", rate.Get())
	// ./binary --help
}
