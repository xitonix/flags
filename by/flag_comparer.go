package by

import (
	"reflect"
	"strings"

	"go.xitonix.io/flags/core"
)

var (
	// DeclarationOrder the default sort order.
	// The flags will be printed in the same order as they have been defined.
	DeclarationOrder Comparer = nil
	// LongNameAscending sort by long name in ascending order.
	LongNameAscending = FlagComparer{Ascending: true, Field: LongName}
	// LongNameDescending sort by long name in descending order.
	LongNameDescending = FlagComparer{Field: LongName}
	// ShortNameAscending sort by short name in ascending order.
	ShortNameAscending = FlagComparer{Ascending: true, Field: ShortName}
	// ShortNameDescending sort by short name in descending order.
	ShortNameDescending = FlagComparer{Field: ShortName}
	// KeyAscending sort by key in ascending order.
	KeyAscending = FlagComparer{Ascending: true, Field: Key}
	// KeyDescending sort by key in descending order.
	KeyDescending = FlagComparer{Field: Key}
	// UsageAscending sort by usage in ascending order.
	UsageAscending = FlagComparer{Ascending: true, Field: Usage}
	// UsageDescending sort by usage in descending order.
	UsageDescending = FlagComparer{Field: Usage}
)

// ComparisonField the field the comparison will be based on.
type ComparisonField int8

const (
	// LongName long name field
	LongName ComparisonField = iota
	// ShortName short name field
	ShortName
	// Key key field
	Key
	//Usage usage field
	Usage
)

// FlagComparer represents an implementation of the by.Comparer interface.
type FlagComparer struct {
	// Ascending applies the comparison function in ascending order.
	Ascending bool
	// Field specifies the field to compare.
	Field ComparisonField
}

// LessThan returns true if the specified filed of f1 is less than the same field in f2.
//
// If Ascending is false, the ordering will be reversed.
func (f FlagComparer) LessThan(f1, f2 core.Flag) bool {
	if reflect.ValueOf(f1).IsNil() || reflect.ValueOf(f2).IsNil() {
		return false
	}
	field1, field2 := f.getFieldsToCompare(f1, f2)
	result := strings.Compare(field1, field2)
	if f.Ascending {
		return result < 0
	}
	return result > 0
}

func (f FlagComparer) getFieldsToCompare(f1, f2 core.Flag) (string, string) {
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
