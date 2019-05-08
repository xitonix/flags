package flags

import "strings"

func isEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
