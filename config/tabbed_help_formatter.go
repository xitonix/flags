package config

import (
	"fmt"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type TabbedHelpFormatter struct {
}

func (*TabbedHelpFormatter) Format(f core.Flag) string {
	if f.IsHidden() {
		return ""
	}
	short := f.ShortName()
	if !internal.IsEmpty(short) {
		short = "-" + short + ","
	}
	var def string
	if dv := f.Default(); dv != nil {
		def = fmt.Sprintf(" (default: %v)", dv)
	}

	return fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s\n", short, f.LongName(), f.Env().Name(), f.Type(), f.Usage(), def)
}
