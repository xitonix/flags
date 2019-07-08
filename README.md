[![Build Status](https://travis-ci.org/xitonix/flags.svg?branch=master)](https://travis-ci.org/xitonix/flags)
[![codecov](https://codecov.io/gh/xitonix/flags/branch/master/graph/badge.svg)](https://codecov.io/gh/xitonix/flags)
# flags

Package flags is a POSIX/GNU compliant flags library providing a simple, yet fully flexible API to manage command line arguments.

The value of each flag can be provided by different sources. Two built-in value providers are command line argument and environment variable sources, with the former at the beginning of the chain, meaning the values parsed by the command line argument source will override the environment variable values. The package also provides the API to register new custom sources to the chain with a desired priority. For example, you may have your own implementation of the Source interface to read from a YAML or JSON file.

The API is packed with a full set of standard built in flag types, from int to IP address and many more. But you can also build a flag for your custom type and ask the library to pass it through the processing pipeline, the same way it treats any pre-built flags.

## Features

- Built-in flag types
  - bool and [ ]bool
  - byte
  - CIDR and [ ]CIDR
  - Counter
  - Duration and [ ]Duration
  - Datetime/Date/Timestamp
  - float32/float64/[ ]float64
  - int/int8/int16/int32/int64/[ ]int
  - IP address and [ ]IP
  - string and [ ]string
  - map[string]string
  - uint/uint8/uint16/uint32/uint64/[ ]uint
  - verbosity
- Ability to add `Hidden`, `Deprecated` and `Required` flags
- Pre-built command line argument and **environment variable** sources
- Automatic key generation
- Flag value validation through callbacks and providing valid lists
- API extendability to read the flag values from custom sources
- Fully customisable help formatter and sort order

## Installation

```bash
go get github.com/xitonix/flags
```



## Usage

```go
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

```

