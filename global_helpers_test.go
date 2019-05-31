package flags_test

import (
	"testing"

	"go.xitonix.io/flags"
	"go.xitonix.io/flags/mocks"
)

func TestEnableAutoKeyGeneration(t *testing.T) {
	flags.EnableAutoKeyGeneration()
	if !flags.DefaultBucket.Options().AutoKeys {
		t.Error("The default bucket's auto key generation was expected to be ON")
	}
}

func TestSetKeyPrefix(t *testing.T) {
	prefix := "prefix"
	expected := "PREFIX"
	flags.SetKeyPrefix(prefix)
	actual := flags.DefaultBucket.Options().KeyPrefix
	if actual != expected {
		t.Errorf("The default bucket's key prefix was expected to be %s, but it was %s", expected, actual)
	}
}

func TestSetLogger(t *testing.T) {
	lg := &mocks.Logger{}
	flags.SetLogger(lg)
	actual := flags.DefaultBucket.Options().Logger
	if actual != lg {
		t.Error("The default bucket's logger has not been set as expected")
	}
}

func TestParse(t *testing.T) {
	flags.DefaultBucket = flags.NewBucket()
	flags.Options().Terminator = &mocks.Terminator{}
	flags.Options().Logger = &mocks.Logger{}
	flags.Options().HelpProvider.Writer = mocks.NewInMemoryWriter()
	flags.String("long", "usage")
	flags.Parse()
	actual := len(flags.DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
}

func TestOptions(t *testing.T) {
	tm := &mocks.Terminator{}
	flags.Options().Terminator = tm
	if flags.DefaultBucket.Options().Terminator != tm {
		t.Error("The default bucket's options has not been set as expected")
	}
}

func TestGlobalString(t *testing.T) {
	flags.DefaultBucket = flags.NewBucket()
	flags.String("long", "usage")
	actual := len(flags.DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
}

func TestGlobalStringP(t *testing.T) {
	flags.DefaultBucket = flags.NewBucket()
	flags.StringP("long", "s", "usage")
	actual := len(flags.DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
}
