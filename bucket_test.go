package flags

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/config"
	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/mocks"
	"github.com/xitonix/flags/test"
)

func TestBucket_Parse_Deprecated_And_Required(t *testing.T) {
	w := mocks.NewInMemoryWriter()
	lg := &mocks.Logger{}
	tm := &mocks.Terminator{}
	env := mocks.NewEnvReader()
	bucket := newBucket([]string{"--flag"}, env,
		config.WithHelpWriter(w),
		config.WithLogger(lg),
		config.WithTerminator(tm))

	f := mocks.NewFlag("flag", "f").Required().MarkAsDeprecated()
	bucket.flags = append(bucket.flags, f)

	bucket.Parse()

	if !tm.IsTerminated {
		t.Errorf("Expectced to terminate, but it did not happen")
	}

	if tm.Code != core.FailureExitCode {
		t.Errorf("Expectced termination code %d, actual: %d", core.FailureExitCode, tm.Code)
	}

	expectedErr := "-f, --flag is marked as deprecated. An obsolete flag cannot be mandatory"
	if !test.ErrorContainsExact(lg.Error, expectedErr) {
		t.Errorf("Expected '%v', but received %v", expectedErr, lg.Error)
	}
}

func TestBucket_Parse_Required(t *testing.T) {
	w := mocks.NewInMemoryWriter()
	lg := &mocks.Logger{}
	tm := &mocks.Terminator{}
	env := mocks.NewEnvReader()
	bucket := newBucket([]string{}, env,
		config.WithHelpWriter(w),
		config.WithLogger(lg),
		config.WithTerminator(tm))

	f := mocks.NewFlag("flag", "f").Required()
	bucket.flags = append(bucket.flags, f)

	bucket.Parse()

	if !tm.IsTerminated {
		t.Errorf("Expectced to terminate, but it did not happen")
	}

	if tm.Code != core.FailureExitCode {
		t.Errorf("Expectced termination code %d, actual: %d", core.FailureExitCode, tm.Code)
	}

	expectedErr := "-f, --flag flag is required."
	if !test.ErrorContainsExact(lg.Error, expectedErr) {
		t.Errorf("Expected '%v', but received %v", expectedErr, lg.Error)
	}
}

func TestBucket_Parse_Validation(t *testing.T) {
	testCases := []struct {
		title                   string
		args                    []string
		flags                   []core.Flag
		expectedErr             string
		mustPrintHelp           bool
		expectedTerminationCode int
		mustTerminate           bool
	}{
		{
			title:                   "unknown flags",
			args:                    []string{"--unexpected"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			expectedErr:             "is an unknown flag",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "invalid short flag",
			args:                    []string{"--unexpected"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "invalid-short-name")},
			expectedErr:             "can only be a single character",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "long name with single dash",
			args:                    []string{"-long"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			expectedErr:             "is an unknown flag",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "short name with double dash",
			args:                    []string{"--f"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			expectedErr:             "is an unknown flag",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "reserved flags",
			args:                    []string{"flag"},
			flags:                   []core.Flag{mocks.NewFlag("help", "h")},
			expectedErr:             "reserved",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "flags with the same long names",
			flags:                   []core.Flag{mocks.NewFlag("flag", "f1"), mocks.NewFlag("flag", "f2")},
			expectedErr:             "already exists",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "flags with the same short names",
			flags:                   []core.Flag{mocks.NewFlag("flag1", "f"), mocks.NewFlag("flag2", "f")},
			expectedErr:             "already exists",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			for _, flag := range tc.flags {
				bucket.flags = append(bucket.flags, flag)
			}

			if bucket.Options() == nil {
				t.Error("Bucket options was not supposed to be nil")
			}

			bucket.Parse()

			if tc.mustTerminate || tm.IsTerminated {
				testTermination(t, tc.mustTerminate, tm.IsTerminated, tc.expectedTerminationCode, tm.Code)
				return
			}

			if tc.mustPrintHelp && w.WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && w.WriteCounter != 0 {
				t.Errorf("Did not expect the Help() function to get called, but it happened")
			}

			if !test.ErrorContains(lg.Error, tc.expectedErr) {
				t.Errorf("Expected '%v', but received %v", tc.expectedErr, lg.Error)
			}
		})
	}
}

func TestBucket_Parse_Help_Request(t *testing.T) {
	testCases := []struct {
		title                   string
		args                    []string
		flags                   []core.Flag
		mustPrintHelp           bool
		expectedTerminationCode int
	}{
		{
			title:                   "help requested with help flag and no other registered flags",
			args:                    []string{"--help"},
			mustPrintHelp:           false,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with help flag and other registered flags",
			args:                    []string{"--help"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag and no other registered flags",
			args:                    []string{"-h"},
			mustPrintHelp:           false,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag and other registered flags",
			args:                    []string{"-h"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag set to true and no other registered flags",
			args:                    []string{"-h=true"},
			mustPrintHelp:           false,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag set to true and other registered flags",
			args:                    []string{"-h=true"},
			flags:                   []core.Flag{mocks.NewFlag("flag", "f")},
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			for _, flag := range tc.flags {
				bucket.flags = append(bucket.flags, flag)
			}

			bucket.Parse()

			if !tm.IsTerminated {
				t.Errorf("Expected to terminate, but it did not happen")
			}

			if tm.Code != tc.expectedTerminationCode {
				t.Errorf("Expected termination code: %d, Actual: %d", tc.expectedTerminationCode, tm.Code)
			}

			if tc.mustPrintHelp && w.WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && w.WriteCounter != 0 {
				t.Errorf("Did not expect the Help() function to get called, but it happened")
			}
		})
	}
}

func TestBucket_Parse_Help_Sort(t *testing.T) {
	testCases := []struct {
		title         string
		args          []string
		flags         []core.Flag
		comparer      by.Comparer
		expectedLines []string
	}{
		{
			title:         "default order as declared xa",
			args:          []string{"--help"},
			comparer:      by.DeclarationOrder,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x"), mocks.NewFlag("a-long", "a")},
			expectedLines: []string{"x-long", "a-long"},
		},
		{
			title:         "default order as declared ax",
			args:          []string{"--help"},
			comparer:      by.DeclarationOrder,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a"), mocks.NewFlag("x-long", "x")},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "sort by long name ascending",
			args:          []string{"--help"},
			comparer:      by.LongNameAscending,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x"), mocks.NewFlag("a-long", "a")},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "sort by long name descending",
			args:          []string{"--help"},
			comparer:      by.LongNameDescending,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a"), mocks.NewFlag("x-long", "x")},
			expectedLines: []string{"x-long", "a-long"},
		},
		{
			title:         "sort by short name ascending",
			args:          []string{"--help"},
			comparer:      by.ShortNameAscending,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x"), mocks.NewFlag("a-long", "a")},
			expectedLines: []string{"a", "x"},
		},
		{
			title:         "sort by short name descending",
			args:          []string{"--help"},
			comparer:      by.ShortNameDescending,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a"), mocks.NewFlag("x-long", "x")},
			expectedLines: []string{"x", "a"},
		},
		{
			title:         "sort by key ascending",
			args:          []string{"--help"},
			comparer:      by.KeyAscending,
			flags:         []core.Flag{mocks.NewFlagWithKey("x-long", "x", "x-key"), mocks.NewFlagWithKey("a-long", "a", "a-key")},
			expectedLines: []string{"A_KEY", "X_KEY"},
		},
		{
			title:         "sort by key descending",
			args:          []string{"--help"},
			comparer:      by.KeyDescending,
			flags:         []core.Flag{mocks.NewFlagWithKey("a-long", "a", "a-key"), mocks.NewFlagWithKey("x-long", "x", "x-key")},
			expectedLines: []string{"X_KEY", "A_KEY"},
		},
		{
			title:         "sort by usage ascending",
			args:          []string{"--help"},
			comparer:      by.UsageAscending,
			flags:         []core.Flag{mocks.NewFlagWithUsage("x-long", "x", "x usage"), mocks.NewFlagWithUsage("a-long", "a", "a usage")},
			expectedLines: []string{"a usage", "x usage"},
		},
		{
			title:         "sort by usage descending",
			args:          []string{"--help"},
			comparer:      by.UsageDescending,
			flags:         []core.Flag{mocks.NewFlagWithUsage("a-long", "a", "a usage"), mocks.NewFlagWithUsage("x-long", "x", "x usage")},
			expectedLines: []string{"x usage", "a usage"},
		},
		{
			title:         "required flags first",
			args:          []string{"--help"},
			comparer:      by.RequiredFirst,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x"), mocks.NewFlag("a-long", "a").Required()},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "required flags last",
			args:          []string{"--help"},
			comparer:      by.RequiredLast,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x").Required(), mocks.NewFlag("a-long", "a")},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "deprecated flags first",
			args:          []string{"--help"},
			comparer:      by.DeprecatedFirst,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x"), mocks.NewFlag("a-long", "a").MarkAsDeprecated()},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "deprecated flags last",
			args:          []string{"--help"},
			comparer:      by.DeprecatedLast,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x").MarkAsDeprecated(), mocks.NewFlag("a-long", "a")},
			expectedLines: []string{"a-long", "x-long"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm),
				config.WithSortOrder(tc.comparer))

			for _, flag := range tc.flags {
				bucket.flags = append(bucket.flags, flag)
			}

			bucket.Parse()

			if !tm.IsTerminated {
				t.Errorf("Expected to terminate, but it did not happen")
			}

			if tm.Code != core.SuccessExitCode {
				t.Errorf("Expected termination code: %d, Actual: %d", core.SuccessExitCode, tm.Code)
			}

			if w.WriteCounter != 2 {
				t.Errorf("Expectced to call Help() for 2 flags, Actual: %d", w.WriteCounter)
			}

			for i, line := range w.Lines {
				if !strings.Contains(line, tc.expectedLines[i]) {
					t.Errorf("Expectced the Help line %d to contain '%s', Actual %s", i+1, tc.expectedLines[i], line)
				}
			}
		})
	}
}

func TestBucket_Parse_Help_Failure(t *testing.T) {
	testCases := []struct {
		title            string
		forceWriteToFail bool
		forceCloseToFail bool
		expectedError    string
	}{
		{
			title:         "no failure",
			expectedError: "",
		},
		{
			title:            "write failing",
			forceWriteToFail: true,
			expectedError:    mocks.ErrExpected.Error(),
		},
		{
			title:            "close failing",
			forceCloseToFail: true,
			expectedError:    mocks.ErrExpected.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			w.ForceCloseToBreak = tc.forceCloseToFail
			w.ForceWriteToBreak = tc.forceWriteToFail

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket([]string{}, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			bucket.flags = []core.Flag{mocks.NewFlag("long", "s")}

			bucket.Help()
			if !test.ErrorContains(lg.Error, tc.expectedError) {
				t.Errorf("Expected to receive '%s' error, but received: %v", mocks.ErrExpected, lg.Error)
			}
		})
	}
}

func TestFlags(t *testing.T) {
	w := mocks.NewInMemoryWriter()
	lg := &mocks.Logger{}
	tm := &mocks.Terminator{}
	env := mocks.NewEnvReader()
	bucket := newBucket([]string{}, env,
		config.WithHelpWriter(w),
		config.WithLogger(lg),
		config.WithTerminator(tm))

	_ = bucket.String("long", "usage")

	if len(bucket.Flags()) != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", len(bucket.flags))
	}
}

func TestBucket_Parse_Value_Args_Source(t *testing.T) {
	testCases := []struct {
		title         string
		expectedValue interface{}
		defaultValue  string
		args          []string
		flag          *mocks.Flag
		makeSetToFail bool
	}{
		{
			title:         "long name provided",
			flag:          mocks.NewFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: "flag_value",
		},
		{
			title:         "short name provided",
			flag:          mocks.NewFlag("flag", "f"),
			args:          []string{"-f", "flag_value"},
			expectedValue: "flag_value",
		},
		{
			title:         "no flag is provided with default value",
			flag:          mocks.NewFlag("flag", "f"),
			defaultValue:  "default",
			expectedValue: "default",
		},
		{
			title:         "no flag is provided without default value",
			flag:          mocks.NewFlag("flag", "f"),
			expectedValue: "",
		},
		{
			title:         "make Set call to fail",
			flag:          mocks.NewFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: nil,
			makeSetToFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.MakeSetToFail = tc.makeSetToFail

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.makeSetToFail {
				if !test.ErrorContains(lg.Error, mocks.ErrExpected.Error()) {
					t.Errorf("Expected to receive '%s' error, but received: %v", mocks.ErrExpected, lg.Error)
				}
				if !tm.IsTerminated {
					t.Errorf("Expected to terminate, but it didn't happen")
				}
				if tm.Code != core.FailureExitCode {
					t.Errorf("Expectced termination code %d, actual: %d", core.FailureExitCode, tm.Code)
				}
			}

			if tc.flag.Get() != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.Get())
			}
		})
	}
}

func TestBucket_Callbacks(t *testing.T) {
	testCases := []struct {
		title              string
		expectedValue      interface{}
		defaultValue       string
		args               []string
		flag               *mocks.Flag
		makeSetToFail      bool
		setPreSetCallback  bool
		preSetCallbackErr  error
		setPostSetCallback bool
		postSetCallbackErr error
	}{
		{
			title:         "no call back has been provided",
			flag:          mocks.NewFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: "flag_value",
		},
		{
			title:             "pre set callback",
			flag:              mocks.NewFlag("flag", "f"),
			args:              []string{"--flag", "flag_value"},
			expectedValue:     "flag_value",
			setPreSetCallback: true,
			preSetCallbackErr: nil,
		},
		{
			title:             "pre set callback failure",
			flag:              mocks.NewFlag("flag", "f"),
			args:              []string{"--flag", "flag_value"},
			expectedValue:     nil,
			setPreSetCallback: true,
			preSetCallbackErr: errors.New("pre-set callback error"),
		},
		{
			title:              "post set callback",
			flag:               mocks.NewFlag("flag", "f"),
			args:               []string{"--flag", "flag_value"},
			expectedValue:      "flag_value",
			setPostSetCallback: true,
			postSetCallbackErr: nil,
		},
		{
			title:              "post set callback failure",
			flag:               mocks.NewFlag("flag", "f"),
			args:               []string{"--flag", "flag_value"},
			expectedValue:      "flag_value",
			setPostSetCallback: true,
			postSetCallbackErr: errors.New("post-set callback error"),
		},
		{
			title:              "pre and post set callback",
			flag:               mocks.NewFlag("flag", "f"),
			args:               []string{"--flag", "flag_value"},
			expectedValue:      "flag_value",
			setPreSetCallback:  true,
			preSetCallbackErr:  nil,
			setPostSetCallback: true,
			postSetCallbackErr: nil,
		},
		{
			title:         "make Set to fail",
			flag:          mocks.NewFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: nil,
			makeSetToFail: true,
		},
		{
			title:              "make Set to fail with post set callback",
			flag:               mocks.NewFlag("flag", "f"),
			args:               []string{"--flag", "flag_value"},
			expectedValue:      nil,
			makeSetToFail:      true,
			setPostSetCallback: true,
			postSetCallbackErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.MakeSetToFail = tc.makeSetToFail

			bucket.flags = []core.Flag{tc.flag}

			var (
				preSetCallbackCalled  bool
				postSetCallbackCalled bool
			)

			if tc.setPreSetCallback {
				bucket.Options().PreSetCallback = func(flag core.Flag, value string) error {
					if value != "flag_value" {
						t.Errorf("Pre-set callback, Expected 'flag_value', Received: %s", value)
					}
					preSetCallbackCalled = true
					return tc.preSetCallbackErr
				}
			}

			if tc.setPostSetCallback {
				bucket.Options().PostSetCallback = func(flag core.Flag, value string) error {
					if value != "flag_value" {
						t.Errorf("Post-set callback, Expected 'flag_value', Received: %s", value)
					}
					postSetCallbackCalled = true
					return tc.postSetCallbackErr
				}
			}

			bucket.Parse()

			if tc.setPreSetCallback && !preSetCallbackCalled {
				t.Error("Pre-set callback has not been called")
			}

			if tc.preSetCallbackErr != nil {
				checkError(t, tc.preSetCallbackErr, lg.Error, tm, "pre-set callback")
				if tc.flag.IsSet() {
					t.Error("The flag value was not supposed to be set, because pre-set callback has failed")
				}
			}

			if tc.makeSetToFail && postSetCallbackCalled {
				t.Error("Post-set callback was not supposed to bet called")
			}

			if tc.setPostSetCallback && !tc.makeSetToFail && !postSetCallbackCalled {
				t.Error("Post-set callback has not been called")
			}

			if tc.postSetCallbackErr != nil {
				checkError(t, tc.postSetCallbackErr, lg.Error, tm, "post-set callback")
			}

			if tc.flag.Get() != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.Get())
			}
		})
	}
}

func checkError(t *testing.T, expected, actual error, tm *mocks.Terminator, operation string) {
	t.Helper()
	if !test.ErrorContains(expected, actual.Error()) {
		t.Errorf("Expected %s error: '%s', but received: %s", operation, expected, actual)
	}

	if !tm.IsTerminated {
		t.Errorf("Expected to terminate due to %s failure, but it didn't happen", operation)
	}
	if tm.Code != core.FailureExitCode {
		t.Errorf("Expectced termination code %d due to %s failure, actual: %d", core.FailureExitCode, operation, tm.Code)
	}
}

func TestBucket_Parse_Chained_Short_Forms(t *testing.T) {
	testCases := []struct {
		title              string
		expectedValue      map[string]interface{}
		defaultValueSuffix string
		args               []string
		flags              []*mocks.Flag
	}{
		{
			title: "chained short forms without value",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g")},
			args: []string{"-fg"},
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "",
			},
		},
		{
			title: "chained short forms with value",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g")},
			args: []string{"-fg", "value"},
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "value",
			},
		},
		{
			title: "chained short forms with value and equal sign",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g")},
			args: []string{"-fg=value"},
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "value",
			},
		},
		{
			title: "chained short forms without value and with default",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g")},
			args:               []string{"-fg"},
			defaultValueSuffix: "default",
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "",
			},
		},
		{
			title: "chained short forms with value and with default",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g")},
			args:               []string{"-gf", "value"},
			defaultValueSuffix: "default",
			expectedValue: map[string]interface{}{
				"flag-1": "value",
				"flag-2": "",
			},
		},
		{
			title: "chained short mixed with long form",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g"),
				mocks.NewFlag("flag-3", "h")},
			args:               []string{"-fg", "--flag-3"},
			defaultValueSuffix: "default",
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "",
				"flag-3": "",
			},
		},
		{
			title: "chained short mixed with long form and values",
			flags: []*mocks.Flag{
				mocks.NewFlag("flag-1", "f"),
				mocks.NewFlag("flag-2", "g"),
				mocks.NewFlag("flag-3", "h")},
			args:               []string{"-fg", "g-value", "--flag-3=value"},
			defaultValueSuffix: "default",
			expectedValue: map[string]interface{}{
				"flag-1": "",
				"flag-2": "g-value",
				"flag-3": "value",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			bucket.flags = make([]core.Flag, len(tc.flags))
			for i, flag := range tc.flags {
				if tc.defaultValueSuffix != "" {
					flag.SetDefaultValue(flag.LongName() + tc.defaultValueSuffix)
				}
				bucket.flags[i] = flag
			}

			bucket.Parse()

			for _, flag := range tc.flags {
				actual := flag.Get()
				if actual != tc.expectedValue[flag.LongName()] {
					t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue[flag.LongName()], actual)
				}
			}
		})
	}
}

func TestBucket_Parse_Value_Environment_Variable_Source(t *testing.T) {
	testCases := []struct {
		title            string
		envVariableValue string
		expectedValue    interface{}
		defaultValue     string
		args             []string
		flag             *mocks.Flag
		keyPrefix        string
		autoKey          bool
		explicitKey      string
		envVariableKey   string
		MakeSetToFail    bool
	}{
		{
			title:         "default value only",
			flag:          mocks.NewFlag("flag", "f"),
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
		{
			title:            "with default value and explicit environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "env_value",
		},
		{
			title:            "with default value and explicit empty environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "",
			expectedValue:    "",
		},
		{
			title:            "without default value and with explicit environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			envVariableValue: "env_value",
			expectedValue:    "env_value",
		},
		{
			title:            "without default value and with explicit empty environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			envVariableValue: "",
			expectedValue:    "",
		},
		{
			title:            "with auto environment variable and no prefix",
			flag:             mocks.NewFlag("test-flag", "f"),
			autoKey:          true,
			envVariableKey:   "TEST_FLAG",
			envVariableValue: "env_value",
			expectedValue:    "env_value",
		},
		{
			title:            "with empty auto environment variable and no prefix",
			flag:             mocks.NewFlag("test-flag", "f"),
			autoKey:          true,
			envVariableKey:   "TEST_FLAG",
			envVariableValue: "",
			expectedValue:    "",
		},
		{
			title:            "with auto environment variable and prefix",
			flag:             mocks.NewFlag("test-flag", "f"),
			autoKey:          true,
			keyPrefix:        "Prefix",
			envVariableKey:   "PREFIX_TEST_FLAG",
			envVariableValue: "env_value",
			expectedValue:    "env_value",
		},
		{
			title:            "with empty auto environment variable and prefix",
			flag:             mocks.NewFlag("test-flag", "f"),
			autoKey:          true,
			keyPrefix:        "Prefix",
			envVariableKey:   "PREFIX_TEST_FLAG",
			envVariableValue: "",
			expectedValue:    "",
		},
		{
			title:            "long flag value override with default value and explicit environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			args:             []string{"--flag", "flag_value"},
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
		{
			title:            "long flag value override with default value and auto environment variable",
			flag:             mocks.NewFlag("test-flag", "f"),
			args:             []string{"--test-flag", "flag_value"},
			autoKey:          true,
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
		{
			title:            "long flag value override with default value and prefixed auto environment variable",
			flag:             mocks.NewFlag("test-flag", "f"),
			args:             []string{"--test-flag", "flag_value"},
			autoKey:          true,
			keyPrefix:        "Prefix",
			envVariableKey:   "PREFIX_TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
		{
			title:            "short flag value override with default value and explicit environment variable",
			flag:             mocks.NewFlag("flag", "f"),
			args:             []string{"-f", "flag_value"},
			explicitKey:      "TEST_FLAG",
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
		{
			title:            "short flag value override with default value and auto environment variable",
			flag:             mocks.NewFlag("test-flag", "f"),
			args:             []string{"-f", "flag_value"},
			autoKey:          true,
			envVariableKey:   "TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
		{
			title:            "short flag value override with default value and prefixed auto environment variable",
			flag:             mocks.NewFlag("test-flag", "f"),
			args:             []string{"-f", "flag_value"},
			autoKey:          true,
			keyPrefix:        "Prefix",
			envVariableKey:   "PREFIX_TEST_FLAG",
			defaultValue:     "default_value",
			envVariableValue: "env_value",
			expectedValue:    "flag_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm),
				config.WithKeyPrefix(tc.keyPrefix))

			bucket.opts.AutoKeys = tc.autoKey

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.MakeSetToFail = tc.MakeSetToFail

			if tc.explicitKey != "" {
				tc.flag.Key().SetID(tc.explicitKey)
			}

			bucket.flags = []core.Flag{tc.flag}

			if tc.envVariableKey != "" {
				env.Set(tc.envVariableKey, tc.envVariableValue)
			}

			bucket.Parse()

			if tc.MakeSetToFail {
				if !test.ErrorContains(lg.Error, mocks.ErrExpected.Error()) {
					t.Errorf("Expected to receive '%s' error, but received: %v", mocks.ErrExpected, lg.Error)
				}
				if !tm.IsTerminated {
					t.Errorf("Expected to terminate, but it didn't happen")
				}
				if tm.Code != core.FailureExitCode {
					t.Errorf("Expectced termination code %d, actual: %d", core.FailureExitCode, tm.Code)
				}
				return
			}

			if tc.flag.Get() != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.Get())
			}
		})
	}
}

func TestBucket_Parse_Custom_Source(t *testing.T) {
	type valueList struct {
		defaultVal string
		cli        string
		env        string
		custom     string
	}
	testCases := []struct {
		title         string
		index         int
		setKey        bool
		values        valueList
		expectedValue string
	}{
		// Custom Source at the end
		{
			title:         "all values provided with custom source at the end",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "cli", env: "env", custom: "custom"},
			expectedValue: "cli",
		},
		{
			title:         "no command line argument with custom source at the end",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "env", custom: "custom"},
			expectedValue: "env",
		},
		{
			title:         "only custom value with custom source at the end",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "only default value with custom source at the end",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: ""},
			expectedValue: "default",
		},
		{
			title:         "no values provided with custom source at the end",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "", cli: "", env: "", custom: ""},
			expectedValue: "",
		},

		// Custom Source in the middle
		{
			title:         "all values provided with custom source in the middle",
			index:         1,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "cli", env: "env", custom: "custom"},
			expectedValue: "cli",
		},
		{
			title:         "no command line argument with custom source in the middle",
			index:         1,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "env", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "only custom value with custom source in the middle",
			index:         1,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "only default value with custom source in the middle",
			index:         1,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: ""},
			expectedValue: "default",
		},
		{
			title:         "no values provided with custom source in the middle",
			index:         1,
			setKey:        true,
			values:        valueList{defaultVal: "", cli: "", env: "", custom: ""},
			expectedValue: "",
		},

		// Custom Source at the beginning
		{
			title:         "all values provided with custom source at the beginning",
			index:         0,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "cli", env: "env", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "no command line argument with custom source at the beginning",
			index:         0,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "env", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "only custom value with custom source at the beginning",
			index:         0,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: "custom"},
			expectedValue: "custom",
		},
		{
			title:         "only default value with custom source at the beginning",
			index:         0,
			setKey:        true,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: ""},
			expectedValue: "default",
		},
		{
			title:         "no values provided with custom source at the beginning",
			index:         0,
			setKey:        true,
			values:        valueList{defaultVal: "", cli: "", env: "", custom: ""},
			expectedValue: "",
		},

		// Without Key
		{
			title:         "all values provided with custom source at the end and no key",
			index:         2,
			setKey:        false,
			values:        valueList{defaultVal: "default", cli: "cli", env: "env", custom: "custom"},
			expectedValue: "cli",
		},
		{
			title:         "no command line argument with custom source at the end and no key",
			index:         2,
			setKey:        false,
			values:        valueList{defaultVal: "default", cli: "", env: "env", custom: "custom"},
			expectedValue: "default",
		},
		{
			title:         "only custom value with custom source at the end and no key",
			index:         2,
			setKey:        false,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: "custom"},
			expectedValue: "default",
		},
		{
			title:         "only default value with custom source at the end and no key",
			index:         2,
			setKey:        false,
			values:        valueList{defaultVal: "default", cli: "", env: "", custom: ""},
			expectedValue: "default",
		},
		{
			title:         "no values provided with custom source at the end and no key",
			index:         2,
			setKey:        true,
			values:        valueList{defaultVal: "", cli: "", env: "", custom: ""},
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			arguments := make([]string, 0)

			if tc.values.cli != "" {
				arguments = append(arguments, "--flag", tc.values.cli)
			}

			flag := mocks.NewFlag("flag", "f")
			if tc.setKey {
				flag = flag.WithKey("FLAG")
			}

			if tc.values.env != "" {
				env.Set(flag.Key().String(), tc.values.env)
			}

			if tc.values.defaultVal != "" {
				flag.SetDefaultValue(tc.values.defaultVal)
			}

			bucket := newBucket(arguments, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			bucket.flags = []core.Flag{flag}

			src := NewMemorySource()
			if tc.values.custom != "" {
				src.Add(flag.Key().String(), tc.values.custom)
			}

			bucket.AddSource(src, tc.index)

			bucket.Parse()
			actual := flag.Get()
			if actual != tc.expectedValue {
				t.Errorf("Expected value: %v, Actual: %v", tc.expectedValue, actual)
			}
		})
	}
}

func TestBucket_KeyGeneration(t *testing.T) {
	testCases := []struct {
		title                string
		keyPrefix            string
		autoKeys             bool
		explicitKey          string
		expectedKeyValue     string
		expectedFlagPrefix   string
		expectedBucketPrefix string
		flag                 core.Flag
	}{
		{
			title:                "prefix without auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "Prefix",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "",
			autoKeys:             false,
		},
		{
			title:                "prefix with auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "Prefix",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_FLAG",
			autoKeys:             true,
		},
		{
			title:                "no prefix without auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "",
			expectedBucketPrefix: "",
			expectedKeyValue:     "",
			autoKeys:             false,
		},
		{
			title:                "no prefix with auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "",
			expectedBucketPrefix: "",
			expectedKeyValue:     "FLAG",
			autoKeys:             true,
		},
		{
			title:                "prefix with explicit key ID",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "Prefix",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_EXPLICIT_KEY",
			autoKeys:             false,
		},
		{
			title:                "not prefixed with explicit key ID",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "",
			expectedKeyValue:     "EXPLICIT_KEY",
			autoKeys:             false,
		},
		{
			title:                "prefix with explicit key ID and auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "Prefix",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_EXPLICIT_KEY",
			autoKeys:             true,
		},
		{
			title:                "not prefixed with explicit key ID and auto generation",
			flag:                 mocks.NewFlag("flag", "f"),
			keyPrefix:            "",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "",
			expectedKeyValue:     "EXPLICIT_KEY",
			autoKeys:             true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket([]string{}, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithKeyPrefix(tc.keyPrefix),
				config.WithTerminator(tm))

			bucket.opts.AutoKeys = tc.autoKeys

			if tc.explicitKey != "" {
				tc.flag.Key().SetID(tc.explicitKey)
			}

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.expectedKeyValue != tc.flag.Key().String() {
				t.Errorf("Expected Key: %v, Actual: %s", tc.expectedKeyValue, tc.flag.Key())
			}

			if tc.expectedBucketPrefix != bucket.opts.KeyPrefix {
				t.Errorf("Expected Bucket Prefix: %v, Actual: %v", tc.expectedBucketPrefix, bucket.opts.KeyPrefix)
			}

			if tc.autoKeys != bucket.opts.AutoKeys {
				t.Errorf("Expected Auto Key Generation: %v, Actual: %v", tc.autoKeys, bucket.opts.AutoKeys)
			}
		})
	}
}

func TestBucket_AddSource(t *testing.T) {
	testCases := []struct {
		title          string
		src            core.Source
		index          int
		expected       map[int]core.Source
		expectedLength int
	}{
		{
			title: "nil source",
			src:   nil,
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
			},
			expectedLength: 2,
		},
		{
			title: "add source to the beginning",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &MemorySource{},
				1: &argSource{},
				2: &envVariableSource{},
			},
			index:          0,
			expectedLength: 3,
		},
		{
			title: "add source to the beginning with negative index",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &MemorySource{},
				1: &argSource{},
				2: &envVariableSource{},
			},
			index:          -1,
			expectedLength: 3,
		},
		{
			title: "add source to the end",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
				2: &MemorySource{},
			},
			index:          2,
			expectedLength: 3,
		},
		{
			title: "add source to the end with out of range index",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
				2: &MemorySource{},
			},
			index:          200,
			expectedLength: 3,
		},
		{
			title: "add source to the middle",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &argSource{},
				1: &MemorySource{},
				2: &envVariableSource{},
			},
			index:          1,
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			if len(tc.expected) == 0 {
				t.Error("The expected source list cannot be empty")
			}

			bucket := NewBucket()
			bucket.AddSource(tc.src, tc.index)

			if len(bucket.sources) != tc.expectedLength {
				t.Errorf("Expected Number of Sources: %d, Actual: %d", tc.expectedLength, len(bucket.sources))
				return
			}

			for i, expected := range tc.expected {
				actual := bucket.sources[i]
				if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
					t.Errorf("Expected Source at index %d: %T, Actual: %T", i, expected, actual)
				}
			}
		})
	}
}

func TestBucket_AppendSource(t *testing.T) {
	testCases := []struct {
		title          string
		src            core.Source
		expected       map[int]core.Source
		expectedLength int
	}{
		{
			title: "nil source",
			src:   nil,
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
			},
			expectedLength: 2,
		},
		{
			title: "non nil source",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
				2: &MemorySource{},
			},
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			if len(tc.expected) == 0 {
				t.Error("The expected source list cannot be empty")
			}

			bucket := NewBucket()
			bucket.AppendSource(tc.src)

			if len(bucket.sources) != tc.expectedLength {
				t.Errorf("Expected Number of Sources: %d, Actual: %d", tc.expectedLength, len(bucket.sources))
				return
			}

			for i, expected := range tc.expected {
				actual := bucket.sources[i]
				if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
					t.Errorf("Expected Source at index %d: %T, Actual: %T", i, expected, actual)
				}
			}
		})
	}
}

func TestBucket_PrependSource(t *testing.T) {
	testCases := []struct {
		title          string
		src            core.Source
		expected       map[int]core.Source
		expectedLength int
	}{
		{
			title: "nil source",
			src:   nil,
			expected: map[int]core.Source{
				0: &argSource{},
				1: &envVariableSource{},
			},
			expectedLength: 2,
		},
		{
			title: "non nil source",
			src:   NewMemorySource(),
			expected: map[int]core.Source{
				0: &MemorySource{},
				1: &argSource{},
				2: &envVariableSource{},
			},
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			if len(tc.expected) == 0 {
				t.Error("The expected source list cannot be empty")
			}

			bucket := NewBucket()
			bucket.PrependSource(tc.src)

			if len(bucket.sources) != tc.expectedLength {
				t.Errorf("Expected Number of Sources: %d, Actual: %d", tc.expectedLength, len(bucket.sources))
				return
			}

			for i, expected := range tc.expected {
				actual := bucket.sources[i]
				if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
					t.Errorf("Expected Source at index %d: %T, Actual: %T", i, expected, actual)
				}
			}
		})
	}
}

func TestBucket_Add(t *testing.T) {
	bucket := NewBucket()
	f := mocks.NewFlag("long", "short")
	bucket.Add(f)
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	af := bucket.Flags()[0]
	if _, ok := af.(*mocks.Flag); !ok {
		t.Errorf("Expected %T, but received %T", &mocks.Flag{}, af)
	}
}

func TestBucket_String(t *testing.T) {
	bucket := NewBucket()
	bucket.String("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringFlag{}, f)
	}
}

func TestBucket_StringP(t *testing.T) {
	bucket := NewBucket()
	bucket.StringP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringFlag{}, f)
	}
}

func TestBucket_StringMap(t *testing.T) {
	bucket := NewBucket()
	bucket.StringMap("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringMapFlag{}, f)
	}
}

func TestBucket_StringMapP(t *testing.T) {
	bucket := NewBucket()
	bucket.StringMapP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringMapFlag{}, f)
	}
}

func TestBucket_StringSliceMap(t *testing.T) {
	bucket := NewBucket()
	bucket.StringSliceMap("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringSliceMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceMapFlag{}, f)
	}
}

func TestBucket_StringSliceMapP(t *testing.T) {
	bucket := NewBucket()
	bucket.StringSliceMapP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringSliceMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceMapFlag{}, f)
	}
}

func TestBucket_Int(t *testing.T) {
	bucket := NewBucket()
	bucket.Int("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntFlag{}, f)
	}
}

func TestBucket_IntP(t *testing.T) {
	bucket := NewBucket()
	bucket.IntP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntFlag{}, f)
	}
}

func TestBucket_Int64(t *testing.T) {
	bucket := NewBucket()
	bucket.Int64("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int64Flag{}, f)
	}
}

func TestBucket_Int64P(t *testing.T) {
	bucket := NewBucket()
	bucket.Int64P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int64Flag{}, f)
	}
}

func TestBucket_Int32(t *testing.T) {
	bucket := NewBucket()
	bucket.Int32("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int32Flag{}, f)
	}
}

func TestBucket_Int32P(t *testing.T) {
	bucket := NewBucket()
	bucket.Int32P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int32Flag{}, f)
	}
}

func TestBucket_Int16(t *testing.T) {
	bucket := NewBucket()
	bucket.Int16("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int16Flag{}, f)
	}
}

func TestBucket_Int16P(t *testing.T) {
	bucket := NewBucket()
	bucket.Int16P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int16Flag{}, f)
	}
}

func TestBucket_Int8(t *testing.T) {
	bucket := NewBucket()
	bucket.Int8("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int8Flag{}, f)
	}
}

func TestBucket_Int8P(t *testing.T) {
	bucket := NewBucket()
	bucket.Int8P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Int8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int8Flag{}, f)
	}
}

func TestBucket_UInt(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UIntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntFlag{}, f)
	}
}

func TestBucket_UIntP(t *testing.T) {
	bucket := NewBucket()
	bucket.UIntP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UIntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntFlag{}, f)
	}
}

func TestBucket_UInt64(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt64("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt64Flag{}, f)
	}
}

func TestBucket_UInt64P(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt64P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt64Flag{}, f)
	}
}

func TestBucket_UInt32(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt32("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt32Flag{}, f)
	}
}

func TestBucket_UInt32P(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt32P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt32Flag{}, f)
	}
}

func TestBucket_UInt16(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt16("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt16Flag{}, f)
	}
}

func TestBucket_UInt16P(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt16P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt16Flag{}, f)
	}
}

func TestBucket_UInt8(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt8("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt8Flag{}, f)
	}
}

func TestBucket_UInt8P(t *testing.T) {
	bucket := NewBucket()
	bucket.UInt8P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UInt8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt8Flag{}, f)
	}
}

func TestBucket_Byte(t *testing.T) {
	bucket := NewBucket()
	bucket.Byte("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*ByteFlag); !ok {
		t.Errorf("Expected %T, but received %T", &ByteFlag{}, f)
	}
}

func TestBucket_Bool(t *testing.T) {
	bucket := NewBucket()
	bucket.Bool("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*BoolFlag); !ok {
		t.Errorf("Expected %T, but received %T", &BoolFlag{}, f)
	}
}

func TestBucket_Parse_Bool(t *testing.T) {
	testCases := []struct {
		title         string
		expectedValue bool
		defaultValue  bool
		setDefault    bool
		args          []string
		flag          *BoolFlag
		key           string
		keyValue      string
	}{
		{
			title:         "long name has been provided",
			flag:          Bool("boolean", "usage"),
			args:          []string{"--boolean"},
			expectedValue: true,
		},
		{
			title:         "short name has been provided",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b"},
			expectedValue: true,
		},
		{
			title:         "both short and long names are provided",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b", "--boolean"},
			expectedValue: true,
		},
		{
			title:         "long name with explicit boolean value has been provided no equal sign",
			flag:          Bool("boolean", "usage"),
			args:          []string{"--boolean", "false"},
			expectedValue: false,
		},
		{
			title:         "long name with explicit boolean value has been provided with equal sign",
			flag:          Bool("boolean", "usage"),
			args:          []string{"--boolean=false"},
			expectedValue: false,
		},
		{
			title:         "long name with explicit integer value has been provided no equal sign",
			flag:          Bool("boolean", "usage"),
			args:          []string{"--boolean", "0"},
			expectedValue: false,
		},
		{
			title:         "long name with explicit integer value has been provided with equal sign",
			flag:          Bool("boolean", "usage"),
			args:          []string{"--boolean=0"},
			expectedValue: false,
		},

		{
			title:         "short name with explicit boolean value has been provided no equal sign",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b", "false"},
			expectedValue: false,
		},
		{
			title:         "short name with explicit boolean value has been provided with equal sign",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b=false"},
			expectedValue: false,
		},
		{
			title:         "short name with explicit integer value has been provided no equal sign",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b", "0"},
			expectedValue: false,
		},
		{
			title:         "short name with explicit integer value has been provided with equal sign",
			flag:          Bool("boolean", "usage").WithShort("b"),
			args:          []string{"-b=0"},
			expectedValue: false,
		},
		{
			title:         "with true default value",
			flag:          Bool("boolean", "usage").WithShort("b"),
			setDefault:    true,
			defaultValue:  true,
			expectedValue: true,
		},
		{
			title:         "with false default value",
			flag:          Bool("boolean", "usage").WithShort("b"),
			setDefault:    true,
			defaultValue:  false,
			expectedValue: false,
		},
		{
			title:         "no default value",
			flag:          Bool("boolean", "usage").WithShort("b"),
			setDefault:    false,
			expectedValue: false,
		},
		{
			title:         "empty key has been provided",
			flag:          Bool("boolean", "usage"),
			key:           "BOOLEAN_KEY",
			expectedValue: true,
		},
		{
			title:         "none empty key with true value has been provided",
			flag:          Bool("boolean", "usage"),
			key:           "BOOLEAN_KEY",
			keyValue:      "true",
			expectedValue: true,
		},
		{
			title:         "none empty key with false value has been provided",
			flag:          Bool("boolean", "usage"),
			key:           "BOOLEAN_KEY",
			keyValue:      "false",
			expectedValue: false,
		},
		{
			title:         "none empty key with true numeric value has been provided",
			flag:          Bool("boolean", "usage"),
			key:           "BOOLEAN_KEY",
			keyValue:      "1",
			expectedValue: true,
		},
		{
			title:         "none empty key with false numeric value has been provided",
			flag:          Bool("boolean", "usage"),
			key:           "BOOLEAN_KEY",
			keyValue:      "0",
			expectedValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			if tc.key != "" {
				tc.flag = tc.flag.WithKey(tc.key)
				env.Set(tc.key, tc.keyValue)
			}
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			if tc.setDefault {
				tc.flag = tc.flag.WithDefault(tc.defaultValue)
			}

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.flag.Get() != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.Get())
			}
		})
	}
}

func TestBucket_BoolSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.BoolSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*BoolSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &BoolSliceFlag{}, f)
	}
}

func TestBucket_Float64(t *testing.T) {
	bucket := NewBucket()
	bucket.Float64("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64Flag{}, f)
	}
}

func TestBucket_Float64P(t *testing.T) {
	bucket := NewBucket()
	bucket.Float64P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64Flag{}, f)
	}
}

func TestBucket_Float32(t *testing.T) {
	bucket := NewBucket()
	bucket.Float32("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float32Flag{}, f)
	}
}

func TestBucket_Float32P(t *testing.T) {
	bucket := NewBucket()
	bucket.Float32P("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float32Flag{}, f)
	}
}

func TestBucket_Counter(t *testing.T) {
	bucket := NewBucket()
	bucket.Counter("long", "usage").WithShort("s")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*CounterFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CounterFlag{}, f)
	}
}

func TestBucket_Parse_Counter(t *testing.T) {
	testCases := []struct {
		title         string
		expectedValue int
		defaultValue  int
		setDefault    bool
		args          []string
		flag          *CounterFlag
		key           string
		keyValue      string
	}{
		{
			title:         "long name has been provided",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter"},
			expectedValue: 1,
		},
		{
			title:         "long name with value has been provided with no equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter", "10"},
			expectedValue: 10,
		},
		{
			title:         "long name with value has been provided with equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter=10"},
			expectedValue: 10,
		},
		{
			title:         "short name has been provided",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-c"},
			expectedValue: 1,
		},
		{
			title:         "short name with value has been provided with no equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-c", "10"},
			expectedValue: 10,
		},
		{
			title:         "short name with value has been provided with equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-c=10"},
			expectedValue: 10,
		},
		{
			title:         "chained short name has been provided",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-cc"},
			expectedValue: 2,
		},
		{
			title:         "short name has been provided multiple times",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-c", "-c", "-c"},
			expectedValue: 3,
		},
		{
			title:         "long name has been provided multiple times",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter", "--counter"},
			expectedValue: 2,
		},
		{
			title:         "long and short names have been provided multiple times",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter", "--counter", "-c", "-c"},
			expectedValue: 4,
		},
		{
			title:         "no argument has been provided",
			flag:          Counter("counter", "usage").WithShort("c"),
			expectedValue: 0,
		},
		{
			title:         "no arguments with default value",
			flag:          Counter("counter", "usage").WithShort("c"),
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 100,
		},
		{
			title:         "long name with default value",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 1,
		},
		{
			title:         "multiple long names with default value",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter", "--counter"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 2,
		},
		{
			title:         "long name with value and default value no equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter", "200"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 200,
		},
		{
			title:         "long name with value and default value with equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"--counter=200"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 200,
		},
		{
			title:         "explicit value of long names overrides the counter value without equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-ccc", "--counter", "1000"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 1000,
		},
		{
			title:         "explicit value of short names overrides the counter value without equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-ccc", "1000"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 1000,
		},
		{
			title:         "explicit value of long names overrides the counter value with equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-ccc", "--counter=1000"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 1000,
		},
		{
			title:         "explicit value of short names overrides the counter value with equal sign",
			flag:          Counter("counter", "usage").WithShort("c"),
			args:          []string{"-ccc=1000"},
			setDefault:    true,
			defaultValue:  100,
			expectedValue: 1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			w := mocks.NewInMemoryWriter()
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			if tc.key != "" {
				tc.flag = tc.flag.WithKey(tc.key)
				env.Set(tc.key, tc.keyValue)
			}
			bucket := newBucket(tc.args, env,
				config.WithHelpWriter(w),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			if tc.setDefault {
				tc.flag = tc.flag.WithDefault(tc.defaultValue)
			}

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.flag.Get() != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.Get())
			}
		})
	}
}

func TestBucket_Verbosity(t *testing.T) {
	bucket := NewBucket()
	bucket.Verbosity("usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*CounterFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CounterFlag{}, f)
	}
	if f.LongName() != "verbose" {
		t.Errorf("Expected Long Name: verbose, Actual %s", f.LongName())
	}
	if f.ShortName() != "v" {
		t.Errorf("Expected Short Name: v, Actual %s", f.ShortName())
	}
}

func TestBucket_Duration(t *testing.T) {
	bucket := NewBucket()
	bucket.Duration("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*DurationFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationFlag{}, f)
	}
}

func TestBucket_DurationP(t *testing.T) {
	bucket := NewBucket()
	bucket.DurationP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*DurationFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationFlag{}, f)
	}
}

func TestBucket_DurationSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.DurationSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*DurationSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationSliceFlag{}, f)
	}
}

func TestBucket_DurationSliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.DurationSliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*DurationSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationSliceFlag{}, f)
	}
}

func TestBucket_Time(t *testing.T) {
	bucket := NewBucket()
	bucket.Time("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*TimeFlag); !ok {
		t.Errorf("Expected %T, but received %T", &TimeFlag{}, f)
	}
}

func TestBucket_TimeP(t *testing.T) {
	bucket := NewBucket()
	bucket.TimeP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*TimeFlag); !ok {
		t.Errorf("Expected %T, but received %T", &TimeFlag{}, f)
	}
}

func TestBucket_StringSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.StringSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceFlag{}, f)
	}
}

func TestBucket_StringSliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.StringSliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*StringSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceFlag{}, f)
	}
}

func TestBucket_IntSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.IntSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntSliceFlag{}, f)
	}
}

func TestBucket_IntSliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.IntSliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntSliceFlag{}, f)
	}
}

func TestBucket_UIntSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.UIntSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UIntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntSliceFlag{}, f)
	}
}

func TestBucket_UIntSliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.UIntSliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*UIntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntSliceFlag{}, f)
	}
}

func TestBucket_Float64Slice(t *testing.T) {
	bucket := NewBucket()
	bucket.Float64Slice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float64SliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64SliceFlag{}, f)
	}
}

func TestBucket_Float64SliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.Float64SliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*Float64SliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64SliceFlag{}, f)
	}
}

func TestBucket_IPAddress(t *testing.T) {
	bucket := NewBucket()
	bucket.IPAddress("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IPAddressFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressFlag{}, f)
	}
}

func TestBucket_IPAddressP(t *testing.T) {
	bucket := NewBucket()
	bucket.IPAddressP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IPAddressFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressFlag{}, f)
	}
}

func TestBucket_IPAddressSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.IPAddressSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IPAddressSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressSliceFlag{}, f)
	}
}

func TestBucket_IPAddressSliceP(t *testing.T) {
	bucket := NewBucket()
	bucket.IPAddressSliceP("long", "s", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*IPAddressSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressSliceFlag{}, f)
	}
}

func TestBucket_CIDR(t *testing.T) {
	bucket := NewBucket()
	bucket.CIDR("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*CIDRFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CIDRFlag{}, f)
	}
}

func TestBucket_CIDRSlice(t *testing.T) {
	bucket := NewBucket()
	bucket.CIDRSlice("long", "usage")
	actual := len(bucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := bucket.Flags()[0]
	if _, ok := f.(*CIDRSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CIDRSliceFlag{}, f)
	}
}

func testTermination(t *testing.T, mustTerminate, isTerminated bool, expectedCode, actualCode int) {
	t.Helper()
	if mustTerminate {
		if !isTerminated {
			t.Errorf("Expectced to terminate, but it did not happen")
		}
		if actualCode != expectedCode {
			t.Errorf("Expectced termination code %d, actual: %d", expectedCode, actualCode)
		}
		return
	}
	if isTerminated {
		t.Errorf("Did not expect to terminate, but the app was terminated")
	}
}
