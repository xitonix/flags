package flags

import (
	"testing"

	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/test"
)

func TestBucket_Parse(t *testing.T) {
	testCases := []struct {
		title                   string
		args                    []string
		flags                   []core.Flag
		expectedErr             string
		mustPrintHelp           bool
		expectedTerminationCode int
		mustTerminate           bool
	}{
		// Validation
		{
			title:                   "unknown flags",
			args:                    []string{"--unexpected"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			expectedErr:             "is an unknown flag",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "reserved flags",
			args:                    []string{"flag"},
			flags:                   []core.Flag{newMockedFlag("help", "h")},
			expectedErr:             "reserved",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "flags with the same long names",
			flags:                   []core.Flag{newMockedFlag("flag", "f1"), newMockedFlag("flag", "f2")},
			expectedErr:             "already exists",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},
		{
			title:                   "flags with the same short names",
			flags:                   []core.Flag{newMockedFlag("flag1", "f"), newMockedFlag("flag2", "f")},
			expectedErr:             "already exists",
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.FailureExitCode,
		},

		// HELP
		{
			title:                   "help requested with help flag and no other registered flags",
			args:                    []string{"--help"},
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: 0,
		},
		{
			title:                   "help requested with help flag and other registered flags",
			args:                    []string{"--help"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag and no other registered flags",
			args:                    []string{"-h"},
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag and other registered flags",
			args:                    []string{"-h"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag set to true and no other registered flags",
			args:                    []string{"-h=true"},
			mustTerminate:           true,
			mustPrintHelp:           false,
			expectedTerminationCode: core.SuccessExitCode,
		},
		{
			title:                   "help requested with H flag set to true and other registered flags",
			args:                    []string{"-h=true"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			mustPrintHelp:           true,
			expectedTerminationCode: core.SuccessExitCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := config.NewHelpProvider(&test.NullWriter{}, &config.TabbedHelpFormatter{})

			lg := &test.LoggerMock{}
			tm := &test.TerminatorMock{}
			bucket := newBucket(tc.args,
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

			if tc.mustPrintHelp && hp.Writer.(*test.NullWriter).WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && hp.Writer.(*test.NullWriter).WriteCounter != 0 {
				t.Errorf("Did not expect the Help() function to get called, but it happened")
			}

			if !test.ErrorContains(lg.Error, tc.expectedErr) {
				t.Errorf("Expected '%v', but received %v", tc.expectedErr, lg.Error)
			}
		})
	}
}

func TestBucket_Parse_Value_Args_Source(t *testing.T) {
	testCases := []struct {
		title         string
		expectedValue string
		defaultValue  string
		args          []string
		flag          *flagMock
		makeSetToFail bool
	}{
		{
			title:         "long name provided",
			flag:          newMockedFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: "flag_value",
		},
		{
			title:         "short name provided",
			flag:          newMockedFlag("flag", "f"),
			args:          []string{"-f", "flag_value"},
			expectedValue: "flag_value",
		},
		{
			title:         "no flag is provided with default value",
			flag:          newMockedFlag("flag", "f"),
			defaultValue:  "default",
			expectedValue: "default",
		},
		{
			title:         "no flag is provided without default value",
			flag:          newMockedFlag("flag", "f"),
			expectedValue: "",
		},
		{
			title:         "make Set call to fail",
			flag:          newMockedFlag("flag", "f"),
			args:          []string{"--flag", "flag_value"},
			expectedValue: "flag_value",
			makeSetToFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := config.NewHelpProvider(&test.NullWriter{}, &config.TabbedHelpFormatter{})

			lg := &test.LoggerMock{}
			tm := &test.TerminatorMock{}
			bucket := newBucket(tc.args,
				config.WithHelpProvider(hp),
				config.WithLogger(lg),
				config.WithTerminator(tm))

			tc.flag.SetDefaultValue(tc.defaultValue)
			tc.flag.makeSetToFail = tc.makeSetToFail

			bucket.flags = []core.Flag{tc.flag}

			bucket.Parse()

			if tc.makeSetToFail {
				if !test.ErrorContains(lg.Error, "asked for it") {
					t.Errorf("Expected to receive 'asked for it' error, but received: %s", lg.Error)
				}
				if !tm.IsTerminated {
					t.Errorf("Expected to terminate, but it didn't happen")
				}
				if tm.Code != core.FailureExitCode {
					t.Errorf("Expectced termination code %d, actual: %d", core.FailureExitCode, tm.Code)
				}
				return
			}

			if tc.flag.value != tc.expectedValue {
				t.Errorf("Expected Value: %v, Actual: %v", tc.expectedValue, tc.flag.value)
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
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "Prefix",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "",
			autoKeys:             false,
		},
		{
			title:                "prefix with auto generation",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "Prefix",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_FLAG",
			autoKeys:             true,
		},
		{
			title:                "no prefix without auto generation",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "",
			expectedBucketPrefix: "",
			expectedKeyValue:     "",
			autoKeys:             false,
		},
		{
			title:                "no prefix with auto generation",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "",
			expectedBucketPrefix: "",
			expectedKeyValue:     "FLAG",
			autoKeys:             true,
		},
		{
			title:                "prefix with explicit key ID",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "Prefix",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_EXPLICIT_KEY",
			autoKeys:             false,
		},
		{
			title:                "not prefixed with explicit key ID",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "",
			expectedKeyValue:     "EXPLICIT_KEY",
			autoKeys:             false,
		},
		{
			title:                "prefix with explicit key ID and auto generation",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "Prefix",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "PREFIX",
			expectedKeyValue:     "PREFIX_EXPLICIT_KEY",
			autoKeys:             true,
		},
		{
			title:                "not prefixed with explicit key ID and auto generation",
			flag:                 newMockedFlag("flag", "f"),
			keyPrefix:            "",
			explicitKey:          "Explicit_Key",
			expectedBucketPrefix: "",
			expectedKeyValue:     "EXPLICIT_KEY",
			autoKeys:             true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := config.NewHelpProvider(&test.NullWriter{}, &config.TabbedHelpFormatter{})

			lg := &test.LoggerMock{}
			tm := &test.TerminatorMock{}
			bucket := newBucket([]string{},
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
