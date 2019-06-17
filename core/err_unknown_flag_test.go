package core_test

import (
	"testing"

	"go.xitonix.io/flags/core"
)

func TestErrUnknownFlag_Error(t *testing.T) {
	err := core.NewUnknownFlagErr("long")
	actual := err.Error()
	expected := "long is an unknown flag"
	if actual != expected {
		t.Errorf("Expected error message: %s, Actual: %s", expected, actual)
	}
}
