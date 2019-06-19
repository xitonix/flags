package flags_test

import (
	"reflect"
	"testing"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/test"
)

func checkSliceFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual, actualVar interface{}) {
	t.Helper()
	if !test.ErrorContains(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkFlagSliceValues(t, expectedValue, actual, actualVar)
}

func checkFlagSliceValues(t *testing.T, expectedValue, actual, actualVar interface{}) {
	t.Helper()
	expected := reflect.TypeOf(expectedValue).Elem()
	if !reflect.DeepEqual(reflect.TypeOf(actual).Elem(), expected) {
		t.Errorf("Expected value: %v, Actual: %v", expectedValue, actual)
	}

	fVar := reflect.ValueOf(actualVar).Elem().Interface()
	if !reflect.DeepEqual(reflect.TypeOf(fVar).Elem(), expected) {
		t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue, fVar)
	}
}

func checkFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual, actualVar interface{}) {
	t.Helper()
	if !test.ErrorContains(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkFlagValues(t, expectedValue, actual, actualVar)
}

func checkFlagValues(t *testing.T, expectedValue, actual, actualVar interface{}) {
	t.Helper()

	if actual != expectedValue {
		t.Errorf("Expected value: %v, Actual: %v", expectedValue, actual)
	}

	fVar := reflect.ValueOf(actualVar).Elem()
	if fVar.Interface() != expectedValue {
		t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue, fVar)
	}
}

func checkFlagInitialState(t *testing.T, f core.Flag, expectedType, expectedUsage, expectedLong, expectedShort string) {
	t.Helper()
	if f.LongName() != expectedLong {
		t.Errorf("Expected Long Name: %s, Actual: %s", expectedLong, f.LongName())
	}
	if f.Usage() != expectedUsage {
		t.Errorf("Expected Usage: %s, Actual: %s", expectedUsage, f.Usage())
	}

	if f.IsDeprecated() {
		t.Error("The flag must not be marked as deprecated by default")
	}

	if f.IsHidden() {
		t.Error("The flag must not be marked as hidden by default")
	}

	if f.IsSet() {
		t.Error("The flag value must not be set initially")
	}

	if f.ShortName() != expectedShort {
		t.Errorf("The short name was expected to be '%s' but it was %s", expectedShort, f.ShortName())
	}

	if f.Default() != nil {
		t.Errorf("The initial default value was expected to be nil, but it was %v", f.Default())
	}

	if f.Type() != expectedType {
		t.Errorf("The flag type was expected to be '%s', but it was %s", expectedType, f.Type())
	}
}
