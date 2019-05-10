package internal

import "strings"

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func SanitiseEnvVarName(value string) string {
	return strings.ToUpper(strings.ReplaceAll(value, "-", "_"))
}
