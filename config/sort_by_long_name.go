package config

import "go.xitonix.io/flags/core"

type ByLongNameAsc []core.Flag

func (s ByLongNameAsc) Len() int {
	return len(s)
}

func (s ByLongNameAsc) Less(i, j int) bool {
	return s[i].LongName() < s[j].LongName()
}

func (s ByLongNameAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
