[![Build Status](https://travis-ci.org/xitonix/flags.svg?branch=master)](https://travis-ci.org/xitonix/flags)
[![codecov](https://codecov.io/gh/xitonix/flags/branch/master/graph/badge.svg)](https://codecov.io/gh/xitonix/flags)
# flags

Package flags is a POSIX/GNU compliant flags library providing a simple, yet fully flexible API to manage command line arguments.

The value of each flag can be provided by different sources. Two built-in value providers are command line argument and environment variable sources, with the former at the beginning of the chain, meaning the values parsed by the command line argument source will override the environment variable values. The package also provides the API to register new custom sources to the chain with a desired priority. For example, you may have your own implementation of the Source interface to read from a YAML or JSON file. See AppendSource, PrependSource and AddSource functions for more details.

The API is packed with a full set of standard built in flag types, from int to IP address and many more. But you can also build a flag for your custom type and ask the library to pass it through the processing pipeline, the same way it treats any pre-built flags.

## Installation

```bash
go get github.com/xitonix/flags
```



## Usage

```go
port := flags.Int("port-number", "Port number").WithShort("p")
log := flags.String("log-file", "The path to the log file").WithDefault("/var/log/service.log").Var()

flags.Parse()
fmt.Println("Port", port.Get())
fmt.Println("Log", *log)
```

