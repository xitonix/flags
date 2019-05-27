package flags

import (
	"testing"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/test"
)

func TestRegistry_Add(t *testing.T) {
	testCases := []struct {
		title string

		first            core.Flag
		firstKey         string
		expectedFirstErr string

		second            core.Flag
		secondKey         string
		expectedSecondErr string
	}{
		{
			title:            "with none empty long name",
			first:            newMockedFlag("long", "short"),
			expectedFirstErr: "",
		},
		{
			title:            "with empty long name",
			first:            newMockedFlag("", "short"),
			expectedFirstErr: core.ErrEmptyFlagName.Error(),
		},
		{
			title:            "with empty short name",
			first:            newMockedFlag("long", ""),
			expectedFirstErr: "",
		},
		{
			title:             "two flags with the same long name",
			first:             newMockedFlag("long", "short-1"),
			expectedFirstErr:  "",
			second:            newMockedFlag("long", "short-2"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:             "two flags with the same short name",
			first:             newMockedFlag("long-1", "short"),
			expectedFirstErr:  "",
			second:            newMockedFlag("long-2", "short"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:             "two flags with the same keys",
			first:             newMockedFlag("long-1", "short-1"),
			firstKey:          "key",
			expectedFirstErr:  "",
			second:            newMockedFlag("long-2", "short-2"),
			secondKey:         "key",
			expectedSecondErr: "flag key already exists",
		},
		{
			title:             "two flags with different keys",
			first:             newMockedFlag("long-1", "short-1"),
			firstKey:          "key-1",
			expectedFirstErr:  "",
			second:            newMockedFlag("long-2", "short-2"),
			secondKey:         "key-2",
			expectedSecondErr: "",
		},
		{
			title:             "two flags with different short name casing",
			first:             newMockedFlag("long-1", "short"),
			expectedFirstErr:  "",
			second:            newMockedFlag("long-2", "SHORT"),
			expectedSecondErr: "",
		},
		{
			title:             "two flags with different long name casing",
			first:             newMockedFlag("long", "short-1"),
			expectedFirstErr:  "",
			second:            newMockedFlag("LONG", "short-2"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:            "reserved help long flag",
			first:            newMockedFlag("help", "short"),
			expectedFirstErr: "reserved",
		},
		{
			title:            "reserved h short flag",
			first:            newMockedFlag("long", "h"),
			expectedFirstErr: "reserved",
		},
		{
			title:            "reserved H short flag",
			first:            newMockedFlag("long", "H"),
			expectedFirstErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			tc.first.Key().Set(tc.firstKey)
			reg := newRegistry()
			err := reg.add(tc.first)
			if !test.ErrorContains(err, tc.expectedFirstErr) {
				t.Errorf("Expected to get '%v' error, but received '%v'", tc.expectedFirstErr, err)
			}

			if tc.second != nil {
				tc.second.Key().Set(tc.secondKey)
				err = reg.add(tc.second)
				if !test.ErrorContains(err, tc.expectedSecondErr) {
					t.Errorf("Second Flag: Expected to get '%v' error, but received '%v'", tc.expectedSecondErr, err)
				}
			}
		})
	}
}
