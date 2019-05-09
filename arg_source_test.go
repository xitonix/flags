package flags

import "testing"

func TestArgSource_Read(t *testing.T) {
	type entry struct {
		key, value string
		ok         bool
	}
	testCases := []struct {
		title         string
		in            []string
		expected      []entry
		expectedCount int
	}{
		{
			title: "nil args must return nothing",
			in:    nil,
			expected: []entry{
				{
					key:   "random",
					value: "",
					ok:    false,
				},
			},
			expectedCount: 0,
		},
		{
			title: "empty args must return nothing",
			in:    make([]string, 0),
			expected: []entry{
				{
					key:   "random",
					value: "",
					ok:    false,
				},
			},
			expectedCount: 0,
		},
		{
			title: "single long form with equal sign",
			in:    []string{"--key=value"},
			expected: []entry{
				{
					key:   "key",
					value: "value",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "single long form without equal sign",
			in:    []string{"--key", "value"},
			expected: []entry{
				{
					key:   "key",
					value: "value",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "single short form with equal sign",
			in:    []string{"-k=value"},
			expected: []entry{
				{
					key:   "k",
					value: "value",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "single short form without equal sign",
			in:    []string{"-k", "value"},
			expected: []entry{
				{
					key:   "k",
					value: "value",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "multiple long form variables with equal sign",
			in:    []string{"--key1=value1", "--key2=value2"},
			expected: []entry{
				{
					key:   "key1",
					value: "value1",
					ok:    true,
				},
				{
					key:   "key2",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "multiple long form variables without equal sign",
			in:    []string{"--key1", "value1", "--key2", "value2"},
			expected: []entry{
				{
					key:   "key1",
					value: "value1",
					ok:    true,
				},
				{
					key:   "key2",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "multiple short form variables with equal sign",
			in:    []string{"-k=value1", "-e=value2"},
			expected: []entry{
				{
					key:   "k",
					value: "value1",
					ok:    true,
				},
				{
					key:   "e",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "multiple short form variables without equal sign",
			in:    []string{"-k", "value1", "-e", "value2"},
			expected: []entry{
				{
					key:   "k",
					value: "value1",
					ok:    true,
				},
				{
					key:   "e",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "multiple long form variables with equal sign must override the same key",
			in:    []string{"--key=value1", "--key=value2"},
			expected: []entry{
				{
					key:   "key",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "multiple long form variables without equal sign must override the same key",
			in:    []string{"--key", "value1", "--key", "value2"},
			expected: []entry{
				{
					key:   "key",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "multiple short form variables with equal sign must override the same key",
			in:    []string{"-k=value1", "-k=value2"},
			expected: []entry{
				{
					key:   "k",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "multiple short form variables without equal sign must override the same key",
			in:    []string{"-k", "value1", "-k", "value2"},
			expected: []entry{
				{
					key:   "k",
					value: "value2",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "multiple long form variables without value",
			in:    []string{"--key1", "--key2"},
			expected: []entry{
				{
					key:   "key1",
					value: "",
					ok:    true,
				},
				{
					key:   "key2",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "multiple short form variables without value",
			in:    []string{"-k", "-e"},
			expected: []entry{
				{
					key:   "k",
					value: "",
					ok:    true,
				},
				{
					key:   "e",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 2,
		},
		{
			title: "mix of multiple long form variables with and without value and with equal sign",
			in:    []string{"--key1", "--key2=value2", "--key3"},
			expected: []entry{
				{
					key:   "key1",
					value: "",
					ok:    true,
				},
				{
					key:   "key2",
					value: "value2",
					ok:    true,
				},
				{
					key:   "key3",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 3,
		},
		{
			title: "mix of multiple long form variables with and without value and without equal sign",
			in:    []string{"--key1", "--key2", "value2", "--key3"},
			expected: []entry{
				{
					key:   "key1",
					value: "",
					ok:    true,
				},
				{
					key:   "key2",
					value: "value2",
					ok:    true,
				},
				{
					key:   "key3",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 3,
		},
		{
			title: "mix of multiple short form variables with and without value and with equal sign",
			in:    []string{"-k", "-e=value2", "-p"},
			expected: []entry{
				{
					key:   "k",
					value: "",
					ok:    true,
				},
				{
					key:   "e",
					value: "value2",
					ok:    true,
				},
				{
					key:   "p",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 3,
		},
		{
			title: "mix of multiple short form variables with and without value and without equal sign",
			in:    []string{"-k", "-e", "value2", "-p"},
			expected: []entry{
				{
					key:   "k",
					value: "",
					ok:    true,
				},
				{
					key:   "e",
					value: "value2",
					ok:    true,
				},
				{
					key:   "p",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 3,
		},
		{
			title: "mix of short and long form variables without value",
			in:    []string{"--long1", "-e", "--long2", "-p"},
			expected: []entry{
				{
					key:   "long1",
					value: "",
					ok:    true,
				},
				{
					key:   "e",
					value: "",
					ok:    true,
				},
				{
					key:   "p",
					value: "",
					ok:    true,
				},
				{
					key:   "long2",
					value: "",
					ok:    true,
				},
			},
			expectedCount: 4,
		},
		{
			title: "mix of short and long form variables with value and with equal sign",
			in:    []string{"--long1=lv1", "-e=s1", "--long2=lv2", "-p=s2"},
			expected: []entry{
				{
					key:   "long1",
					value: "lv1",
					ok:    true,
				},
				{
					key:   "e",
					value: "s1",
					ok:    true,
				},
				{
					key:   "p",
					value: "s2",
					ok:    true,
				},
				{
					key:   "long2",
					value: "lv2",
					ok:    true,
				},
			},
			expectedCount: 4,
		},
		{
			title: "mix of short and long form variables with value and without equal sign",
			in:    []string{"--long1", "lv1", "-e", "s1", "--long2", "lv2", "-p", "s2"},
			expected: []entry{
				{
					key:   "long1",
					value: "lv1",
					ok:    true,
				},
				{
					key:   "e",
					value: "s1",
					ok:    true,
				},
				{
					key:   "p",
					value: "s2",
					ok:    true,
				},
				{
					key:   "long2",
					value: "lv2",
					ok:    true,
				},
			},
			expectedCount: 4,
		},
		{
			title: "mix of short and long form variables with value and mixed equal sign",
			in:    []string{"--long1", "lv1", "-e=s1", "--long2=lv2", "-p", "s2"},
			expected: []entry{
				{
					key:   "long1",
					value: "lv1",
					ok:    true,
				},
				{
					key:   "e",
					value: "s1",
					ok:    true,
				},
				{
					key:   "p",
					value: "s2",
					ok:    true,
				},
				{
					key:   "long2",
					value: "lv2",
					ok:    true,
				},
			},
			expectedCount: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			src, _ := newArgSource(tc.in)
			for _, e := range tc.expected {
				actual, ok := src.Read(e.key)
				if ok != e.ok {
					t.Errorf("Exists (Key: %s), Expected: %v, Actual: %v", e.key, e.ok, ok)
				}
				if actual != e.value {
					t.Errorf("Value, Expected: %v, Actual: %v", e.value, actual)
				}
				if tc.expectedCount != len(src.arguments) {
					t.Errorf("Count, Expected: %v, Actual: %v", tc.expectedCount, len(src.arguments))
				}
			}
		})
	}
}

func TestArgSource_Read_With_Special_Values(t *testing.T) {
	type entry struct {
		key, value string
		ok         bool
	}
	testCases := []struct {
		title         string
		in            []string
		expected      []entry
		expectedCount int
	}{
		{
			title: "long form ptr with hyphened value",
			in:    []string{"--key1=--a=10 --b=20"},
			expected: []entry{
				{
					key:   "key1",
					value: "--a=10 --b=20",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "long form ptr with comma separated value and equal sign",
			in:    []string{"--key1=a,b,c"},
			expected: []entry{
				{
					key:   "key1",
					value: "a,b,c",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "long form ptr with comma separated value and without equal sign",
			in:    []string{"--key1", "a,b,c"},
			expected: []entry{
				{
					key:   "key1",
					value: "a,b,c",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "long form ptr with square braces and equal sign",
			in:    []string{"--key1=[a,b,c]"},
			expected: []entry{
				{
					key:   "key1",
					value: "[a,b,c]",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "long form ptr with square braces and without equal sign",
			in:    []string{"--key1", "[a,b,c]"},
			expected: []entry{
				{
					key:   "key1",
					value: "[a,b,c]",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "short form ptr with hyphened value",
			in:    []string{"-k=--a=10 --b=20"},
			expected: []entry{
				{
					key:   "k",
					value: "--a=10 --b=20",
					ok:    true,
				},
			},
			expectedCount: 1,
		},

		{
			title: "short form ptr with comma separated value and equal sign",
			in:    []string{"-k=a,b,c"},
			expected: []entry{
				{
					key:   "k",
					value: "a,b,c",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "short form ptr with comma separated value and without equal sign",
			in:    []string{"-k", "a,b,c"},
			expected: []entry{
				{
					key:   "k",
					value: "a,b,c",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "short form ptr with square braces and equal sign",
			in:    []string{"-k=[a,b,c]"},
			expected: []entry{
				{
					key:   "k",
					value: "[a,b,c]",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
		{
			title: "short form ptr with square braces and without equal sign",
			in:    []string{"-k", "[a,b,c]"},
			expected: []entry{
				{
					key:   "k",
					value: "[a,b,c]",
					ok:    true,
				},
			},
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			src, _ := newArgSource(tc.in)
			for _, e := range tc.expected {
				actual, ok := src.Read(e.key)
				if ok != e.ok {
					t.Errorf("Exists (Key: %s), Expected: %v, Actual: %v", e.key, e.ok, ok)
				}
				if actual != e.value {
					t.Errorf("Value, Expected: %v, Actual: %v", e.value, actual)
				}
				if tc.expectedCount != len(src.arguments) {
					t.Errorf("Count, Expected: %v, Actual: %v", tc.expectedCount, len(src.arguments))
				}
			}
		})
	}
}

func TestHelpFlags(t *testing.T) {
	testCases := []struct {
		title    string
		in       []string
		expected bool
	}{
		{
			title:    "nil input",
			expected: false,
		},
		{
			title:    "empty input",
			in:       make([]string, 0),
			expected: false,
		},
		{
			title:    "no help flag",
			in:       []string{"--random"},
			expected: false,
		},
		{
			title:    "with H flag",
			in:       []string{"-h"},
			expected: true,
		},
		{
			title:    "with HELP flag",
			in:       []string{"--help"},
			expected: true,
		},
		{
			title:    "with long H flag",
			in:       []string{"--h"},
			expected: false,
		},
		{
			title:    "with HELP flag and equal sign",
			in:       []string{"--help="},
			expected: true,
		},
		{
			title:    "with H flag and equal sign",
			in:       []string{"-h="},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, actual := newArgSource(tc.in)
			if actual != tc.expected {
				t.Errorf("Help Request, Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}
