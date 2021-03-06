package by_test

import (
	"testing"

	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/mocks"
)

func TestStringComparer_LessThan(t *testing.T) {
	testCases := []struct {
		field            by.StringComparisonField
		isAscending      bool
		expectedLessThan bool
		f1, f2           *mocks.Flag
		title            string
	}{
		{
			title:            "nil first flag",
			f1:               nil,
			f2:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "nil second flag",
			f2:               nil,
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "nil flags",
			f2:               nil,
			f1:               nil,
			expectedLessThan: false,
		},
		// Long name
		{
			title:            "equal long name ascending",
			field:            by.LongName,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "f1 long name less than f2 long name ascending",
			field:            by.LongName,
			isAscending:      true,
			f1:               mocks.NewFlag("a-long", "a"),
			f2:               mocks.NewFlag("x-long", "x"),
			expectedLessThan: true,
		},
		{
			title:            "f1 long name greater than f2 long name ascending",
			field:            by.LongName,
			isAscending:      true,
			f1:               mocks.NewFlag("x-long", "x"),
			f2:               mocks.NewFlag("a-long", "a"),
			expectedLessThan: false,
		},
		{
			title:            "equal long name descending",
			field:            by.LongName,
			isAscending:      false,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "f1 long name less than f2 long name descending",
			field:            by.LongName,
			isAscending:      false,
			f1:               mocks.NewFlag("a-long", "a"),
			f2:               mocks.NewFlag("x-long", "x"),
			expectedLessThan: false,
		},
		{
			title:            "f1 long name greater than f2 long name descending",
			field:            by.LongName,
			isAscending:      false,
			f1:               mocks.NewFlag("x-long", "x"),
			f2:               mocks.NewFlag("a-long", "a"),
			expectedLessThan: true,
		},
		// Short name
		{
			title:            "equal short name ascending",
			field:            by.ShortName,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "f1 short name less than f2 short name ascending",
			field:            by.ShortName,
			isAscending:      true,
			f1:               mocks.NewFlag("a-long", "a"),
			f2:               mocks.NewFlag("x-long", "x"),
			expectedLessThan: true,
		},
		{
			title:            "f1 short name greater than f2 short name ascending",
			field:            by.ShortName,
			isAscending:      true,
			f1:               mocks.NewFlag("x-long", "x"),
			f2:               mocks.NewFlag("a-long", "a"),
			expectedLessThan: false,
		},

		{
			title:            "equal short name descending",
			field:            by.ShortName,
			isAscending:      false,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "f1 short name less than f2 short name descending",
			field:            by.ShortName,
			isAscending:      false,
			f1:               mocks.NewFlag("a-long", "a"),
			f2:               mocks.NewFlag("x-long", "x"),
			expectedLessThan: false,
		},
		{
			title:            "f1 short name greater than f2 short name descending",
			field:            by.ShortName,
			isAscending:      false,
			f1:               mocks.NewFlag("x-long", "x"),
			f2:               mocks.NewFlag("a-long", "a"),
			expectedLessThan: true,
		},

		// Key
		{
			title:            "equal key ascending",
			field:            by.Key,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s").WithKey("k"),
			f1:               mocks.NewFlag("long", "s").WithKey("k"),
			expectedLessThan: false,
		},
		{
			title:            "f1 key less than f2 key ascending",
			field:            by.Key,
			isAscending:      true,
			f1:               mocks.NewFlag("a-long", "a").WithKey("a"),
			f2:               mocks.NewFlag("x-long", "x").WithKey("x"),
			expectedLessThan: true,
		},
		{
			title:            "f1 key greater than f2 key ascending",
			field:            by.Key,
			isAscending:      true,
			f1:               mocks.NewFlag("x-long", "x").WithKey("x"),
			f2:               mocks.NewFlag("a-long", "a").WithKey("a"),
			expectedLessThan: false,
		},
		{
			title:            "equal key descending",
			field:            by.Key,
			isAscending:      false,
			f2:               mocks.NewFlag("long", "s").WithKey("k"),
			f1:               mocks.NewFlag("long", "s").WithKey("k"),
			expectedLessThan: false,
		},
		{
			title:            "f1 key less than f2 key descending",
			field:            by.Key,
			isAscending:      false,
			f1:               mocks.NewFlag("a-long", "a").WithKey("a"),
			f2:               mocks.NewFlag("x-long", "x").WithKey("x"),
			expectedLessThan: false,
		},
		{
			title:            "f1 key greater than f2 key descending",
			field:            by.Key,
			isAscending:      false,
			f1:               mocks.NewFlag("x-long", "x").WithKey("x"),
			f2:               mocks.NewFlag("a-long", "a").WithKey("a"),
			expectedLessThan: true,
		},

		// Usage
		{
			title:            "equal usage ascending",
			field:            by.Usage,
			isAscending:      true,
			f2:               mocks.NewFlagWithUsage("long", "s", "u"),
			f1:               mocks.NewFlagWithUsage("long", "s", "u"),
			expectedLessThan: false,
		},
		{
			title:            "f1 usage less than f2 usage ascending",
			field:            by.Usage,
			isAscending:      true,
			f1:               mocks.NewFlagWithUsage("a-long", "a", "a"),
			f2:               mocks.NewFlagWithUsage("x-long", "x", "x"),
			expectedLessThan: true,
		},
		{
			title:            "f1 usage greater than f2 usage ascending",
			field:            by.Usage,
			isAscending:      true,
			f1:               mocks.NewFlagWithUsage("x-long", "x", "x"),
			f2:               mocks.NewFlagWithUsage("a-long", "a", "a"),
			expectedLessThan: false,
		},
		{
			title:            "equal usage descending",
			field:            by.Usage,
			isAscending:      false,
			f2:               mocks.NewFlagWithUsage("long", "s", "u"),
			f1:               mocks.NewFlagWithUsage("long", "s", "u"),
			expectedLessThan: false,
		},
		{
			title:            "f1 usage less than f2 usage descending",
			field:            by.Usage,
			isAscending:      false,
			f1:               mocks.NewFlagWithUsage("a-long", "a", "a"),
			f2:               mocks.NewFlagWithUsage("x-long", "x", "x"),
			expectedLessThan: false,
		},
		{
			title:            "f1 usage greater than f2 usage descending",
			field:            by.Usage,
			isAscending:      false,
			f1:               mocks.NewFlagWithUsage("x-long", "x", "x"),
			f2:               mocks.NewFlagWithUsage("a-long", "a", "a"),
			expectedLessThan: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			c := by.StringComparer{
				Ascending: tc.isAscending,
				Field:     tc.field,
			}
			actual := c.LessThan(tc.f1, tc.f2)
			if actual != tc.expectedLessThan {
				t.Errorf("LessThan() Expected: %v, Actual: %v", tc.expectedLessThan, actual)
			}
		})
	}
}
