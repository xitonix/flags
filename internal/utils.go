package internal

import "strings"

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func SanitiseEnvVarName(value string) string {
	return strings.ToUpper(strings.ReplaceAll(value, "-", "_"))
}

func SanitiseLongName(name string) string {
	return SanitiseShortName(strings.ToLower(name))
}

func SanitiseShortName(name string) string {
	return strings.TrimSpace(strings.TrimLeft(name, "-"))
}
