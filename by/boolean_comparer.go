package by

import (
	"reflect"

	"go.xitonix.io/flags/core"
)

// BooleanComparisonField the boolean field the comparison will be based on.
type BooleanComparisonField int8

const (
	// IsRequired the flag's IsRequired status.
	IsRequired BooleanComparisonField = iota
	// IsDeprecated the flag's deprecation status.
	IsDeprecated
)

// BooleanComparer represents an implementation of the by.Comparer interface to compare boolean values.
type BooleanComparer struct {
	// Ascending applies the comparison function in ascending order.
	Ascending bool
	// Field field to compare.
	Field BooleanComparisonField
}

// LessThan returns true if the f1's specified boolean property is true (in ascending mode).
// In descending mode the method will return true, if the f2's specified property is true.
func (b BooleanComparer) LessThan(f1, f2 core.Flag) bool {
	if reflect.ValueOf(f1).IsNil() || reflect.ValueOf(f2).IsNil() {
		return false
	}
	if b.Ascending {
		return b.getBooleanFieldToCompare(f1)
	}
	return b.getBooleanFieldToCompare(f2)
}

func (b BooleanComparer) getBooleanFieldToCompare(f core.Flag) bool {
	switch b.Field {
	case IsRequired:
		return f.IsRequired()
	case IsDeprecated:
		return f.IsDeprecated()
	default:
		return false
	}
}
