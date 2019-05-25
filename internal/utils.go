package internal

import (
	"regexp"
	"strings"
)

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func SanitiseFlagID(value string) string {
	duplicates := regexp.MustCompile(`\s+|[-]+`)
	value = duplicates.ReplaceAllString(strings.TrimSpace(value), "_")
	return strings.ToUpper(value)
}

func SanitiseLongName(name string) string {
	return SanitiseShortName(strings.ToLower(name))
}

func SanitiseShortName(name string) string {
	space := regexp.MustCompile(`\s+`)
	name = space.ReplaceAllString(strings.TrimSpace(name), "-")
	return strings.TrimLeft(name, "-")
}
