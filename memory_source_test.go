package flags_test

import (
	"testing"

	"github.com/xitonix/flags"
)

func TestMemorySource_Add(t *testing.T) {
	testCases := []struct {
		title         string
		key, value    string
		expectedOk    bool
		expectedValue string
	}{
		{
			title:         "empty key",
			key:           "",
			value:         "value",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			title:         "white space key",
			key:           "    ",
			value:         "value",
			expectedValue: "",
			expectedOk:    false,
		},
		{
			title:         "empty value",
			key:           "key",
			value:         "",
			expectedValue: "",
			expectedOk:    true,
		},
		{
			title:         "white space value",
			key:           "key",
			value:         "    ",
			expectedValue: "    ",
			expectedOk:    true,
		},
		{
			title:         "non empty key and value",
			key:           "key",
			value:         "value",
			expectedValue: "value",
			expectedOk:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			m := flags.NewMemorySource()
			m.Add(tc.key, tc.value)
			actual, ok := m.Read(tc.key)
			if ok != tc.expectedOk || actual != tc.expectedValue {
				t.Errorf("Expected Value for Key %s: '%s', Actual: '%s'", tc.key, tc.expectedValue, actual)
			}
		})
	}
}

func TestMemorySource_AddRange(t *testing.T) {
	testCases := []struct {
		title    string
		kv       map[string]string
		expected map[string]struct {
			ok    bool
			value string
		}
	}{
		{
			title: "empty key",
			kv:    map[string]string{"": "value"},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"": {
					ok:    false,
					value: "",
				},
			},
		},
		{
			title: "white space key",
			kv:    map[string]string{"   ": "value"},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"   ": {
					ok:    false,
					value: "",
				},
			},
		},
		{
			title: "empty value",
			kv:    map[string]string{"key": ""},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"key": {
					ok:    true,
					value: "",
				},
			},
		},
		{
			title: "white space value",
			kv:    map[string]string{"key": "  "},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"key": {
					ok:    true,
					value: "  ",
				},
			},
		},
		{
			title: "non empty key and value",
			kv:    map[string]string{"key": "value"},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"key": {
					ok:    true,
					value: "value",
				},
			},
		},
		{
			title: "multiple keys and values",
			kv:    map[string]string{"key1": "value1", "key2": "value2"},
			expected: map[string]struct {
				ok    bool
				value string
			}{
				"key1": {
					ok:    true,
					value: "value1",
				},
				"key2": {
					ok:    true,
					value: "value2",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			m := flags.NewMemorySource()
			m.AddRange(tc.kv)
			for key, expected := range tc.expected {
				actual, ok := m.Read(key)
				if ok != expected.ok || actual != expected.value {
					t.Errorf("Expected Value for Key %s: '%s', Actual: '%s'", key, expected.value, actual)
				}
			}
		})
	}
}
