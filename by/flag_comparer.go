package by

import (
	"strings"

	"go.xitonix.io/flags/core"
)

var (
	Declared            Comparer = nil
	LongNameAscending            = flagComparer{ascending: true, field: longName}
	LongNameDescending           = flagComparer{field: longName}
	ShortNameAscending           = flagComparer{ascending: true, field: shortName}
	ShortNameDescending          = flagComparer{field: shortName}
	KeyAscending                 = flagComparer{ascending: true, field: key}
	KeyDescending                = flagComparer{field: key}
	UsageAscending               = flagComparer{ascending: true, field: usage}
	UsageDescending              = flagComparer{field: usage}
)

type compareField int8

const (
	longName compareField = iota
	shortName
	key
	usage
)

type flagComparer struct {
	ascending bool
	field     compareField
}

func (l flagComparer) LessThan(f1, f2 core.Flag) bool {
	if f1 == nil || f2 == nil {
		return false
	}
	field1, field2 := l.getFieldsToCompare(f1, f2)
	result := strings.Compare(field1, field2)
	if l.ascending {
		return result < 0
	}
	return result > 0
}

func (l flagComparer) getFieldsToCompare(f1, f2 core.Flag) (string, string) {
	switch l.field {
	case longName:
		return f1.LongName(), f2.LongName()
	case shortName:
		return f1.ShortName(), f2.ShortName()
	case key:
		return f1.Key().Get(), f2.Key().Get()
	default:
		return f1.Usage(), f2.Usage()
	}

}
