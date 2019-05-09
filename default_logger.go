package flags

import (
	"fmt"
	"os"
)

type DefaultLogger struct{}

func (*DefaultLogger) Fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
