package by_test

import (
	"testing"

	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/mocks"
)

func TestBooleanComparer_LessThan(t *testing.T) {
	testCases := []struct {
		field            by.BooleanComparisonField
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
		// IsRequired
		{
			title:            "required f1 ascending",
			field:            by.IsRequired,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").Required(),
			expectedLessThan: true,
		},
		{
			title:            "required f1 descending",
			isAscending:      false,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").Required(),
			expectedLessThan: false,
		},
		{
			title:            "required f2 ascending",
			isAscending:      true,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s").Required(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "required f2 descending",
			isAscending:      false,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s").Required(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: true,
		},

		{
			title:            "both required ascending",
			isAscending:      true,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s").Required(),
			f1:               mocks.NewFlag("long", "s").Required(),
			expectedLessThan: true,
		},
		{
			title:            "both required descending",
			isAscending:      false,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s").Required(),
			f1:               mocks.NewFlag("long", "s").Required(),
			expectedLessThan: true,
		},

		{
			title:            "none required ascending",
			isAscending:      true,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "none required descending",
			isAscending:      false,
			field:            by.IsRequired,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},

		// IsDeprecated
		{
			title:            "deprecated f1 ascending",
			field:            by.IsDeprecated,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: true,
		},
		{
			title:            "deprecated f1 descending",
			isAscending:      false,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: false,
		},
		{
			title:            "deprecated f2 ascending",
			isAscending:      true,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "deprecated f2 descending",
			isAscending:      false,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: true,
		},

		{
			title:            "both deprecated ascending",
			isAscending:      true,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: true,
		},
		{
			title:            "both deprecated descending",
			isAscending:      false,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: true,
		},

		{
			title:            "none deprecated ascending",
			isAscending:      true,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "none deprecated descending",
			isAscending:      false,
			field:            by.IsDeprecated,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},

		// IsDeprecated
		{
			title:            "invalid field with deprecated f1 ascending",
			field:            3,
			isAscending:      true,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: false,
		},
		{
			title:            "invalid field with deprecated f1 descending",
			isAscending:      false,
			field:            3,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: false,
		},
		{
			title:            "invalid field with deprecated f2 ascending",
			isAscending:      true,
			field:            3,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "invalid field with deprecated f2 descending",
			isAscending:      false,
			field:            3,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},

		{
			title:            "invalid field with both deprecated ascending",
			isAscending:      true,
			field:            3,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: false,
		},
		{
			title:            "invalid field with both deprecated descending",
			isAscending:      false,
			field:            3,
			f2:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			f1:               mocks.NewFlag("long", "s").MarkAsDeprecated(),
			expectedLessThan: false,
		},

		{
			title:            "invalid field with none of the flags deprecated ascending",
			isAscending:      true,
			field:            3,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
		{
			title:            "invalid field with none of the flags deprecated descending",
			isAscending:      false,
			field:            3,
			f2:               mocks.NewFlag("long", "s"),
			f1:               mocks.NewFlag("long", "s"),
			expectedLessThan: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			c := by.BooleanComparer{
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
