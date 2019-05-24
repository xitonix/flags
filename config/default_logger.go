package config

import (
	"fmt"
)

type DefaultLogger struct{}

func (*DefaultLogger) Print(msg string) {
	fmt.Println(msg)
}
