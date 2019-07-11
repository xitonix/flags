![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/tag-pre/xitonix/flags.svg?label=pre-release)
[![GoDoc](https://godoc.org/github.com/xitonix/flags?status.svg)](https://godoc.org/github.com/xitonix/flags)
[![Go Report Card](https://goreportcard.com/badge/github.com/xitonix/flags)](https://goreportcard.com/report/github.com/xitonix/flags)
[![Build Status](https://travis-ci.org/xitonix/flags.svg?branch=master)](https://travis-ci.org/xitonix/flags)
[![codecov](https://codecov.io/gh/xitonix/flags/branch/master/graph/badge.svg)](https://codecov.io/gh/xitonix/flags)
# Flags

Package flags is a POSIX/GNU compliant flags library providing a simple, yet fully flexible API to manage command line arguments.

The value of each flag can be provided by different sources. Two built-in value providers are **command line argument** and **environment variable** sources, with the former at the beginning of the chain, meaning the values parsed by the command line argument source will override the environment variable values. The package also provides the API to add new custom sources to the chain with a desired priority. For example, you may have your own implementation of the `Source` interface to read from a YAML or JSON file.

The API is packed with a full set of standard built-in flag types, from `int` to IP address and many more. But you can also build a flag for your custom types and ask the library to pass them through the processing pipeline, the same way it treats the pre-built flags.

#### NOTE
**The `pre-release` API may still be subject to breaking changes.**

## Features

- Fully tested with 3K+ unit tests
- Built-in flag types
  - bool and []bool
  - byte
  - CIDR and []CIDR
  - Counter
  - Duration and []Duration
  - Datetime/Date/Timestamp
  - float32/float64/[]float64
  - int/int8/int16/int32/int64/[]int
  - IP address and []IP
  - string and []string
  - map[string]string
  - uint/uint8/uint16/uint32/uint64/[]uint
  - verbosity
  
- Ability to mark the flags as `Hidden`, `Deprecated` and `Required`

- Pre-built command line argument and environment variable sources

- Automatic key generation (For environment variables and other custom sources)

- Flag value validation through callbacks and providing valid lists

- API extendability to read the flag values from custom sources

- Fully customisable help formatter

- Built-in predicates to control the order in which the flags will be printed in the help output

  - Long name (ASC/DESC)
  - Short name (ASC/DESC)
  - Key (ASC/DESC)
  - Usage (ASC/DESC)
  - Sort by Required
  - Sort by Deprecated

- Ability to register your own `Comparer` to gain full control over sorting behaviour

- Supporting the following command line argument formats
  

  ```bash
  # Boolean flags
  
  --bool  --bool=true  --bool=false  --bool=1  --bool=0
  -b  -b=true  -b=false  -b=1  -b=0
  
  # Numeric flags (Integers or floating point numbers)
  
  --num=[+/-]10
  --num [+/-]10
  -n=[+/-]10
  -n [+/-]10
  -n[+/-]10
  
  # Mixed short forms

  -n[+/-]10m         # result: n=+/-10, m=0
  -n[+/-]10m[+/-]20  # result: n=+/-10, m=+/-20
  -n[+/-]10m [+/-]20  # result: n=+/-10, m=+/-20
  
  # Non numeric flags
  
  --key="value"  --key "value"  --key value
  -k="value"  -k "value"  -k value
  ```

  

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
  
   // You can optionally set a prefix for all the automatically generated 
   // or explicitly defined flag keys.
   flags.SetKeyPrefix("PFX")

   // Customising the indicator for deprecated and required flags.
   // The markers will be used in the help output to draw the users' attention
   flags.SetDeprecationMark("**DEPRECATED**")
   flags.SetRequiredFlagMark("**")

   // Changing the sort order.
   // You can use pre-built sort predicates or pass your own Comparer to change the
   // order in which the flags will be printed in the help output.
   flags.SetSortOrder(by.LongNameDescending)

   // Each flag must have a mandatory long form (ie. --file-path) and 
   // an OPTIONAL short form (i.e. -F)
   port := flags.Int("port-number", "Port number").WithShort("p")

   // You can ask the package to set the flag to the specified default value whenever
   // it's not explicitly provided by any available Sources.
   // You may also use the Var() function with each flag to access the pointer to 
   // the underlying variable instead of calling the Get() function to read the falg value.
   log := flags.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()

   // You have full control over the acceptable values for each flag. You can either provide
   // a list of allowed values using WithValidRange method:
   weekend := flags.StringSlice("weekends", "Weekends").WithValidRange(true, "Sat, Sun").WithTrimming()
  
   // or set a fully customisable validator using WithValidationCallback method:
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

   // Deprecated flags will be marked in the help output 
   // using a customisable indicator to draw user's attention
   _ = flags.Int("port", "Legacy port number. Use -p, --port-number instead").MarkAsDeprecated()

   // Required flags will be marked in the help output 
   // using a customisable indicator to draw user's attention.
   // The user must explicitly provide value for a required flag.
   // Setting the default value of required flags will have zero effect.
   rate := flags.Float64("rate", "Conversion rate").Required()

   // Hidden flags will not be displayed in the help output.
   hidden := flags.Bool("enable-logging", "Secret flag").Hide()

   // You can explicitly define the key for each flag. These keys will override
   // their automatically generated counterparts.
   t := flags.Time("start-time", "Start time").WithKey("START")
   ttl := flags.Duration("ttl", "Time to live")

   // The value of Counter flags can be increased by repeating the short or the long form
   // for example -cc --counter --counter will the the counter flag to 4.
   counter := flags.Counter("counter", "Repeat counter")
  
   // Verbosity is an alias for Counter("verbose", "usage").WithShort("v")
   verbose := flags.Verbosity("Verbosity. Repeat -v for higher verbosity levels. Example -vv")

   // This callback function will be called before the flag value is being set by a source.
   preCallback := func(flag core.Flag, value string) error {
	  fmt.Printf("%s will be set to %s\n", flag.LongName(), value)
	  return nil
   }

   // This callback function will be called after the flag value has been set by a source.
   // postCallback will not get called if preCallback returns an error.
   postCallback := func(flag core.Flag, value string) error {
	  fmt.Printf("%s has been set to %s\n", flag.LongName(), value)
	  return nil
   }
  
   flags.SetPreSetCallback(preCallback)
   flags.SetPostSetCallback(postCallback)
   flags.Parse()
  
   // You can read the flag value by calling the Get() method
   fmt.Println("Port", port.Get())
   
   // or accessing the underlying pointer if Var() is used when creating the flag
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

  // It's possible to navigate through all the registered flags
  for _, flag := range flags.DefaultBucket.Flags() {
 	fmt.Printf("--%s (%s) %s\n", flag.LongName(), flag.Type(), flag.Usage())
  }
}

```

