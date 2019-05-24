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
		expectedParseErr        string
		expectedAddErr          string
		mustPrintHelp           bool
		expectedTerminationCode int
		mustTerminate           bool
	}{
		{
			title:                   "unknown flags",
			args:                    []string{"--unexpected"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			expectedParseErr:        "is an unknown flag",
			mustPrintHelp:           true,
			mustTerminate:           true,
			expectedTerminationCode: -1,
		},
		{
			title:          "reserved flags",
			args:           []string{"help"},
			flags:          []core.Flag{newMockedFlag("help", "h")},
			expectedAddErr: "reserved",
		},
		{
			title:                   "help requested with help flag and no other registered flags",
			args:                    []string{"--help"},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           false,
		},
		{
			title:                   "help requested with help flag and other registered flags",
			args:                    []string{"--help"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           true,
		},
		{
			title:                   "help requested with H flag and no other registered flags",
			args:                    []string{"-h"},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           false,
		},
		{
			title:                   "help requested with H flag and other registered flags",
			args:                    []string{"-h"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           true,
		},
		{
			title:                   "help requested with H flag set to true and no other registered flags",
			args:                    []string{"-h=true"},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           false,
		},
		{
			title:                   "help requested with H flag set to true and other registered flags",
			args:                    []string{"-h=true"},
			flags:                   []core.Flag{newMockedFlag("flag", "f")},
			mustTerminate:           true,
			expectedTerminationCode: 0,
			mustPrintHelp:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			hp := config.NewHelpProvider(&test.NullWriter{}, &config.TabbedHelpFormatter{})

			lg := &test.LoggerMock{}
			tm := &test.TerminatorMock{}
			bucket := newBucket(tc.args, config.WithHelpProvider(hp),
				config.WithLogger(lg), config.WithTerminator(tm))

			for _, flag := range tc.flags {
				bucket.addFlag(flag)
				if !test.ErrorContains(lg.Error, tc.expectedAddErr) {
					t.Errorf("Expected '%v', but received %v", tc.expectedAddErr, lg.Error)
				}
				if tc.expectedAddErr != "" {
					return
				}
			}

			bucket.Parse()

			if tc.mustPrintHelp && hp.Writer.(*test.NullWriter).WriteCounter == 0 {
				t.Errorf("Expectced the Help() function to get called, but it did not happen")
			}

			if !tc.mustPrintHelp && hp.Writer.(*test.NullWriter).WriteCounter != 0 {
				t.Errorf("Did not expect the Help() function to get called, but it happened")
			}

			if tc.mustTerminate || tm.IsTerminated {
				testTermination(t, tc.mustTerminate, tm.IsTerminated, tc.expectedTerminationCode, tm.Code)
				return
			}

			if !test.ErrorContains(lg.Error, tc.expectedParseErr) {
				t.Errorf("Expected '%v', but received %v", tc.expectedParseErr, lg.Error)
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
		t.Errorf("Did expectc to terminate, but the app was terminated")
	}
}
