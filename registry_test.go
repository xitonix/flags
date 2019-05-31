package flags

import (
	"testing"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/mocks"
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
			first:            mocks.NewFlag("long", "s"),
			expectedFirstErr: "",
		},
		{
			title:            "with an invalid short name",
			first:            mocks.NewFlag("long", "short"),
			expectedFirstErr: "can only be a single character",
		},
		{
			title:            "with empty long name",
			first:            mocks.NewFlag("", "s"),
			expectedFirstErr: core.ErrEmptyFlagName.Error(),
		},
		{
			title:            "with empty short name",
			first:            mocks.NewFlag("long", ""),
			expectedFirstErr: "",
		},
		{
			title:             "two flags with the same long name",
			first:             mocks.NewFlag("long", "s"),
			expectedFirstErr:  "",
			second:            mocks.NewFlag("long", "t"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:             "two flags with the same short name",
			first:             mocks.NewFlag("long-1", "s"),
			expectedFirstErr:  "",
			second:            mocks.NewFlag("long-2", "s"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:             "two flags with the same keys",
			first:             mocks.NewFlag("long-1", "s"),
			firstKey:          "key",
			expectedFirstErr:  "",
			second:            mocks.NewFlag("long-2", "t"),
			secondKey:         "key",
			expectedSecondErr: "flag key already exists",
		},
		{
			title:             "two flags with different keys",
			first:             mocks.NewFlag("long-1", "s"),
			firstKey:          "key-1",
			expectedFirstErr:  "",
			second:            mocks.NewFlag("long-2", "t"),
			secondKey:         "key-2",
			expectedSecondErr: "",
		},
		{
			title:             "two flags with different short name casing",
			first:             mocks.NewFlag("long-1", "s"),
			expectedFirstErr:  "",
			second:            mocks.NewFlag("long-2", "S"),
			expectedSecondErr: "",
		},
		{
			title:             "two flags with different long name casing",
			first:             mocks.NewFlag("long", "s"),
			expectedFirstErr:  "",
			second:            mocks.NewFlag("LONG", "t"),
			expectedSecondErr: "flag already exists",
		},
		{
			title:            "reserved help long flag",
			first:            mocks.NewFlag("help", "s"),
			expectedFirstErr: "reserved",
		},
		{
			title:            "reserved h short flag",
			first:            mocks.NewFlag("long", "h"),
			expectedFirstErr: "reserved",
		},
		{
			title:            "reserved H short flag",
			first:            mocks.NewFlag("long", "H"),
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
