package by

import (
	"strings"

	"go.xitonix.io/flags/core"
)

var (
	DeclarationOrder    Comparer = nil
	LongNameAscending            = FlagComparer{Ascending: true, Field: LongName}
	LongNameDescending           = FlagComparer{Field: LongName}
	ShortNameAscending           = FlagComparer{Ascending: true, Field: ShortName}
	ShortNameDescending          = FlagComparer{Field: ShortName}
	KeyAscending                 = FlagComparer{Ascending: true, Field: Key}
	KeyDescending                = FlagComparer{Field: Key}
	UsageAscending               = FlagComparer{Ascending: true, Field: Usage}
	UsageDescending              = FlagComparer{Field: Usage}
)

type ComparisonField int8

const (
	LongName ComparisonField = iota
	ShortName
	Key
	Usage
)

type FlagComparer struct {
	Ascending bool
	Field     ComparisonField
}

func (f FlagComparer) LessThan(f1, f2 core.Flag) bool {
	if f1 == nil || f2 == nil {
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
		return f1.Key().Get(), f2.Key().Get()
	default:
		return f1.Usage(), f2.Usage()
	}
}
