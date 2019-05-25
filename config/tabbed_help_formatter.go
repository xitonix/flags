package config

import (
	"fmt"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type TabbedHelpFormatter struct {
	DefaultValueFormatString string
	DeprecatedFlagIndicator  string
}

func NewTabbedHelpFormatter(defaultValueFormatString string, deprecatedIndicator string) *TabbedHelpFormatter {
	return &TabbedHelpFormatter{DefaultValueFormatString: defaultValueFormatString, DeprecatedFlagIndicator: deprecatedIndicator}
}

func (t *TabbedHelpFormatter) Format(f core.Flag) string {
	if f.IsHidden() {
		return ""
	}
	short := f.ShortName()
	if !internal.IsEmpty(short) {
		short = "-" + short + ","
	}
	var def string
	if dv := f.Default(); dv != nil {
		def = fmt.Sprintf(" "+t.DefaultValueFormatString, dv)
	}

	var dep string
	if f.IsDeprecated() {
		dep = " " + t.DeprecatedFlagIndicator
	}

	return fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", short, f.LongName(), f.Key().Value(), f.Type(), f.Usage(), def, dep)
}
