package core

import (
	"fmt"

	"go.xitonix.io/flags/internal"
)

// TabbedHelpFormatter represents a tab separated help formatter.
type TabbedHelpFormatter struct{}

// Format returns a tab separated help string for the flag.
func (t *TabbedHelpFormatter) Format(f Flag, deprecationMark, defaultValueFormatString string) string {
	if f.IsHidden() {
		return ""
	}
	short := f.ShortName()
	if !internal.IsEmpty(short) {
		short = "-" + short + ","
	}
	var def string
	if dv := f.Default(); dv != nil && !internal.IsEmpty(defaultValueFormatString) {
		def = fmt.Sprintf(" "+defaultValueFormatString, dv)
	}

	var dep string
	if f.IsDeprecated() && !internal.IsEmpty(deprecationMark) {
		dep = " " + deprecationMark
	}

	return fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", short, f.LongName(), f.Key().Get(), f.Type(), f.Usage(), def, dep)
}
