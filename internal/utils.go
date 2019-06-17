package internal

import (
	"fmt"
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

func GetExpectedValueString(entry interface{}, i, total int) string {
	switch i {
	case total - 2:
		return fmt.Sprintf("%v and ", entry)
	case total - 1:
		return fmt.Sprintf("%v", entry)
	default:
		return fmt.Sprintf("%v, ", entry)
	}
}

func OutOfRangeErr(value interface{}, longName string, valid string, validCount int) error {
	if IsEmpty(valid) {
		return fmt.Errorf("%v is not an acceptable value for --%s.", value, longName)
	}
	plural := " is"
	if validCount > 1 {
		plural = "s are"
	}
	return fmt.Errorf("%v is not an acceptable value for --%s. The expected value%s %s.", value, longName, plural, valid)
}

func InvalidValueErr(value interface{}, longName, flagType string) error {
	return fmt.Errorf("'%v' is not a valid %s value for --%s", value, flagType, longName)
}
