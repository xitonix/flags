package by

import (
	"reflect"
	"strings"

	"github.com/xitonix/flags/core"
)

// StringComparisonField the string field the comparison will be based on.
type StringComparisonField int8

const (
	// LongName long name field
	LongName StringComparisonField = iota
	// ShortName short name field
	ShortName
	// Key key field
	Key
	// Usage usage field
	Usage
)

// StringComparer represents an implementation of the by.Comparer interface.
type StringComparer struct {
	// Ascending applies the comparison function in ascending order.
	Ascending bool
	// Field specifies the field to compare.
	Field StringComparisonField
}

// LessThan returns true if the specified filed of f1 is less than the same field in f2.
//
// If Ascending is false, the ordering will be reversed.
func (f StringComparer) LessThan(f1, f2 core.Flag) bool {
	if reflect.ValueOf(f1).IsNil() || reflect.ValueOf(f2).IsNil() {
		return false
	}
	field1, field2 := f.getStringFieldsToCompare(f1, f2)
	result := strings.Compare(field1, field2)
	if f.Ascending {
		return result < 0
	}
	return result > 0
}

func (f StringComparer) getStringFieldsToCompare(f1, f2 core.Flag) (string, string) {
	switch f.Field {
	case LongName:
		return f1.LongName(), f2.LongName()
	case ShortName:
		return f1.ShortName(), f2.ShortName()
	case Key:
		return f1.Key().String(), f2.Key().String()
	default:
		return f1.Usage(), f2.Usage()
	}
}
