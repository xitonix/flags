package defaults

import (
	"fmt"
)

type Logger struct{}

func (*Logger) Print(err error) {
	fmt.Println(err.Error())
}
