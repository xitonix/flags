package core_test

import (
	"fmt"
	"testing"

	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/mocks"
)

func TestTabbedHelpFormatter_Format(t *testing.T) {
	testCases := []struct {
		title                    string
		defaultValueFormatString string
		deprecatedFormatString   string
		requiredMark             string
		isRequired               bool
		flag                     *mocks.Flag
		isHidden                 bool
		defaultValue             string
		setDefault               bool
		isDeprecated             bool
		expected                 string
	}{
		{
			title:    "hidden flag",
			expected: "",
			isHidden: true,
			flag:     mocks.NewFlag("long", "s"),
		},
		{
			title:    "with long name and usage only",
			expected: fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", "", ""),
			flag:     mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:    "with long and short names and empty usage",
			expected: fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "-s,", "long", "", "generic", "", "", ""),
			flag:     mocks.NewFlagWithUsage("long", "s", ""),
		},
		{
			title:    "with long and short names along with usage",
			expected: fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "-s,", "long", "", "generic", "usage", "", ""),
			flag:     mocks.NewFlagWithUsage("long", "s", "usage"),
		},
		{
			title:                    "with long name and default value along with usage",
			expected:                 fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", " [DEFAULT: default]", ""),
			defaultValue:             "default",
			setDefault:               true,
			defaultValueFormatString: "[DEFAULT: %v]",
			flag:                     mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:        "without default value indicator",
			expected:     fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", "", ""),
			defaultValue: "default",
			setDefault:   true,
			flag:         mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:                  "deprecated with long name and usage",
			expected:               fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", "", " DEP"),
			deprecatedFormatString: "DEP",
			isDeprecated:           true,
			flag:                   mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:        "without deprecated indicator",
			expected:     fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", "", ""),
			isDeprecated: true,
			flag:         mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:      "required without a mark",
			expected:   fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic", "usage", "", ""),
			isRequired: true,
			flag:       mocks.NewFlagWithUsage("long", "", "usage"),
		},
		{
			title:        "required with a mark",
			expected:     fmt.Sprintf("%s\t--%s\t%s\t%s\t\t\t%s%s%s\n", "", "long", "", "generic*", "usage", "", ""),
			requiredMark: "*",
			isRequired:   true,
			flag:         mocks.NewFlagWithUsage("long", "", "usage"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			f := core.TabbedHelpFormatter{}
			if tc.isDeprecated {
				tc.flag = tc.flag.MarkAsDeprecated()
			}

			tc.flag.SetHidden(tc.isHidden)
			if tc.setDefault {
				tc.flag.SetDefaultValue(tc.defaultValue)
			}
			if tc.isRequired {
				tc.flag = tc.flag.Required()
			}
			actual := f.Format(tc.flag, tc.deprecatedFormatString, tc.defaultValueFormatString, tc.requiredMark)
			if actual != tc.expected {
				t.Errorf("Expected formatted result: '%s', Actual: %s", tc.expected, actual)
			}
		})
	}
}
