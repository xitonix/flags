package by

import "go.xitonix.io/flags/core"

type Comparer interface {
	LessThan(f1, f2 core.Flag) bool
}
