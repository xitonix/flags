/*
Package flags is a POSIX/GNU compliant flags library providing a simple, yet fully flexible API to manage command
line arguments.

The value of each flag can be provided by different sources. Two built-in value providers are command line argument and
environment variable sources, with the former at the beginning of the chain, meaning the values parsed by the command line
argument source will override the environment variable values. The package also provides the API to register new custom sources
to the chain with a desired priority. For example, you may have your own implementation of the Source interface to read from a YAML
or JSON file. See AppendSource, PrependSource and AddSource functions for more details.

The API is packed with a full set of standard built in flag types, from int to IP address and many more. But you can also build a
flag for your custom type and ask the library to pass it through the processing pipeline, the same way it treats any pre-built flags.

Usage

	import "github.com/xitonix/flags"

The package uses the concept of Buckets to organise the flags. You may create a new bucket to register your flags, or go with the default
bucket instead.

Use a new bucket instance

	bucket := flags.NewBucket()
	config := bucket.String("config-file", "The path to the configuration file")
	server := bucket.IPAddress("remote-server", "The remote server to connect to").WithShort("r")
	bucket.Parse()

Use the default bucket

	config := flags.String("config-file", "The path to the configuration file")
	server := flags.IPAddressP("remote-server", "The remote server to connect to").WithShort("r")
	flags.Parse()

The following command line argument formats are supported by the library:

Boolean flags

	--bool
	--bool=true
	--bool=false
	--bool=1
	--bool=0
	-b
	-b=true
	-b=false
	-b=1
	-b=0

Numeric flags (Integers or floating point numbers)

	--num=10
	--num 10
	-n=10
	-n 10
	-n10

	// Mixed short forms
	-n10b
	-n10m20
	-n10m 20

Non numeric flags

	--key="value"
	--key "value"
	--key value
	-k="value"
	-k "value"
	-k value

*/
package flags
