package flags

import (
	"strings"
	"testing"

	"go.xitonix.io/flags/by"
	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/mocks"
	"go.xitonix.io/flags/test"
)

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
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			for _, flag := range tc.flags {
				bucket.flags = append(bucket.flags, flag)
			}

			bucket.Parse()

			if tc.mustTerminate || tm.IsTerminated {
				testTermination(t, tc.mustTerminate, tm.IsTerminated, tc.expectedTerminationCode, tm.Code)
				return
			}

			if tc.mustPrintHelp && hp.Writer.(*mocks.InMemoryWriter).WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && hp.Writer.(*mocks.InMemoryWriter).WriteCounter != 0 {
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
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpProvider(hp),
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

			if tc.mustPrintHelp && hp.Writer.(*mocks.InMemoryWriter).WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && hp.Writer.(*mocks.InMemoryWriter).WriteCounter != 0 {
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
			flags:         []core.Flag{mocks.NewFlag("x-long", "x-short"), mocks.NewFlag("a-long", "a-short")},
			expectedLines: []string{"x-long", "a-long"},
		},
		{
			title:         "default order as declared ax",
			args:          []string{"--help"},
			comparer:      by.DeclarationOrder,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a-short"), mocks.NewFlag("x-long", "x-short")},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "sort by long name ascending",
			args:          []string{"--help"},
			comparer:      by.LongNameAscending,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x-short"), mocks.NewFlag("a-long", "a-short")},
			expectedLines: []string{"a-long", "x-long"},
		},
		{
			title:         "sort by long name descending",
			args:          []string{"--help"},
			comparer:      by.LongNameDescending,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a-short"), mocks.NewFlag("x-long", "x-short")},
			expectedLines: []string{"x-long", "a-long"},
		},
		{
			title:         "sort by short name ascending",
			args:          []string{"--help"},
			comparer:      by.ShortNameAscending,
			flags:         []core.Flag{mocks.NewFlag("x-long", "x-short"), mocks.NewFlag("a-long", "a-short")},
			expectedLines: []string{"a-short", "x-short"},
		},
		{
			title:         "sort by short name descending",
			args:          []string{"--help"},
			comparer:      by.ShortNameDescending,
			flags:         []core.Flag{mocks.NewFlag("a-long", "a-short"), mocks.NewFlag("x-long", "x-short")},
			expectedLines: []string{"x-short", "a-short"},
		},
		{
			title:         "sort by key ascending",
			args:          []string{"--help"},
			comparer:      by.KeyAscending,
			flags:         []core.Flag{mocks.NewFlagWithKey("x-long", "x-short", "x-key"), mocks.NewFlagWithKey("a-long", "a-short", "a-key")},
			expectedLines: []string{"A_KEY", "X_KEY"},
		},
		{
			title:         "sort by key descending",
			args:          []string{"--help"},
			comparer:      by.KeyDescending,
			flags:         []core.Flag{mocks.NewFlagWithKey("a-long", "a-short", "a-key"), mocks.NewFlagWithKey("x-long", "x-short", "x-key")},
			expectedLines: []string{"X_KEY", "A_KEY"},
		},
		{
			title:         "sort by usage ascending",
			args:          []string{"--help"},
			comparer:      by.UsageAscending,
			flags:         []core.Flag{mocks.NewFlagWithUsage("x-long", "x-short", "x usage"), mocks.NewFlagWithUsage("a-long", "a-short", "a usage")},
			expectedLines: []string{"a usage", "x usage"},
		},
		{
			title:         "sort by usage descending",
			args:          []string{"--help"},
			comparer:      by.UsageDescending,
			flags:         []core.Flag{mocks.NewFlagWithUsage("a-long", "a-short", "a usage"), mocks.NewFlagWithUsage("x-long", "x-short", "x usage")},
			expectedLines: []string{"x usage", "a usage"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm),
				config.WithSort(tc.comparer))

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

			writer := hp.Writer.(*mocks.InMemoryWriter)

			if writer.WriteCounter != 2 {
				t.Errorf("Expectced to call Help() for 2 flags, Actual: %d", writer.WriteCounter)
			}

			for i, line := range writer.Lines {
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
			hp := core.NewHelpProvider(w, &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket([]string{}, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			bucket.flags = []core.Flag{mocks.NewFlag("long", "short")}

			bucket.Help()
			if !test.ErrorContains(lg.Error, tc.expectedError) {
				t.Errorf("Expected to receive '%s' error, but received: %v", mocks.ErrExpected, lg.Error)
			}
		})
	}
}

func TestFlags(t *testing.T) {
	w := mocks.NewInMemoryWriter()
	hp := core.NewHelpProvider(w, &core.TabbedHelpFormatter{})
	lg := &mocks.Logger{}
	tm := &mocks.Terminator{}
	env := mocks.NewEnvReader()
	bucket := newBucket([]string{}, env,
		config.WithHelpProvider(hp),
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
		MakeSetToFail bool
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
			expectedValue: "flag_value",
			MakeSetToFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.MakeSetToFail = tc.MakeSetToFail

			bucket.flags = []core.Flag{tc.flag}

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
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})
			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket(tc.args, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm),
				config.WithKeyPrefix(tc.keyPrefix))

			bucket.Options().AutoKeys = tc.autoKey

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.MakeSetToFail = tc.MakeSetToFail

			if tc.explicitKey != "" {
				tc.flag.Key().Set(tc.explicitKey)
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
			hp := core.NewHelpProvider(mocks.NewInMemoryWriter(), &core.TabbedHelpFormatter{})

			lg := &mocks.Logger{}
			tm := &mocks.Terminator{}
			env := mocks.NewEnvReader()
			bucket := newBucket([]string{}, env,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithKeyPrefix(tc.keyPrefix),
				config.WithTerminator(tm))

			bucket.opts.AutoKeys = tc.autoKeys

			if tc.explicitKey != "" {
				tc.flag.Key().Set(tc.explicitKey)
			}

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.expectedKeyValue != tc.flag.Key().Get() {
				t.Errorf("Expected Key: %v, Actual: %v", tc.expectedKeyValue, tc.flag.Key().Get())
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
