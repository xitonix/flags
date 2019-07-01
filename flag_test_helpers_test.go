package flags_test

import (
	"net"
	"reflect"
	"testing"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/test"
)

func checkSliceFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual, actualVar interface{}) {
	t.Helper()
	if !test.ErrorContainsExact(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkSliceFlagValues(t, expectedValue, actual, actualVar)
}

func checkSliceFlagValues(t *testing.T, expectedValue, actual, actualVar interface{}) {
	t.Helper()
	expected := reflect.ValueOf(expectedValue).Interface()
	if !reflect.DeepEqual(reflect.ValueOf(actual).Interface(), expected) {
		t.Errorf("Expected value: %v, Actual: %v", expected, actual)
	}

	fVar := reflect.ValueOf(actualVar).Elem().Interface()
	if !reflect.DeepEqual(fVar, expected) {
		t.Errorf("Expected flag variable: %v, Actual: %v", expected, fVar)
	}
}

func checkIPSliceFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual []net.IP, actualVar *[]net.IP) {
	t.Helper()
	if !test.ErrorContainsExact(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkIPSliceFlagValues(t, expectedValue, actual, actualVar)
}

func checkIPSliceFlagValues(t *testing.T, expectedValue, actual []net.IP, actualVar *[]net.IP) {
	t.Helper()
	if len(actual) != len(expectedValue) {
		t.Errorf("Expected value length: %v, Actual length: %v", len(expectedValue), len(actual))
		return
	}
	if len(*actualVar) != len(expectedValue) {
		t.Errorf("Expected variable length: %v, Actual variable length: %v", len(expectedValue), len(actual))
		return
	}
	for i, act := range actual {
		if !expectedValue[i].Equal(act) {
			t.Errorf("Expected value: %v, Actual: %v", expectedValue[i], act)
		}
		av := (*actualVar)[i]
		if !expectedValue[i].Equal(av) {
			t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue[i], av)
		}
	}
}

func checkCIDRSliceFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual []core.CIDR, actualVar *[]core.CIDR) {
	t.Helper()
	if !test.ErrorContainsExact(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkCIDRSliceFlagValues(t, expectedValue, actual, actualVar)
}

func checkCIDRSliceFlagValues(t *testing.T, expectedValue, actual []core.CIDR, actualVar *[]core.CIDR) {
	t.Helper()
	if len(actual) != len(expectedValue) {
		t.Errorf("Expected value length: %v, Actual length: %v", len(expectedValue), len(actual))
		return
	}
	if len(*actualVar) != len(expectedValue) {
		t.Errorf("Expected variable length: %v, Actual variable length: %v", len(expectedValue), len(actual))
		return
	}
	for i, act := range actual {
		if !expectedValue[i].Equal(act) {
			t.Errorf("Expected value: %v, Actual: %v", expectedValue[i], act)
		}
		av := (*actualVar)[i]
		if !expectedValue[i].Equal(av) {
			t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue[i], av)
		}
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

func checkIPFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual net.IP, actualVar *net.IP) {
	t.Helper()
	if !test.ErrorContains(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkIPFlagValues(t, expectedValue, actual, actualVar)
}

func checkIPFlagValues(t *testing.T, expectedValue, actual net.IP, actualVar *net.IP) {
	t.Helper()

	if !actual.Equal(expectedValue) {
		t.Errorf("Expected value: %v, Actual: %v", expectedValue, actual)
	}

	if !(*actualVar).Equal(expectedValue) {
		t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue, *actualVar)
	}
}

func checkCIDRFlag(t *testing.T, f core.Flag, err error, expectedErr string, expectedValue, actual core.CIDR, actualVar *core.CIDR) {
	t.Helper()
	if !test.ErrorContains(err, expectedErr) {
		t.Errorf("Expected to receive an error with '%s', but received %s", expectedErr, err)
	}

	if expectedErr == "" && !f.IsSet() {
		t.Error("IsSet(), Expected value: true, Actual: false")
	}

	checkCIDRFlagValues(t, expectedValue, actual, actualVar)
}

func checkCIDRFlagValues(t *testing.T, expectedValue, actual core.CIDR, actualVar *core.CIDR) {
	t.Helper()

	if !actual.Equal(expectedValue) {
		t.Errorf("Expected value: %v, Actual: %v", expectedValue, actual)
	}

	if !(*actualVar).Equal(expectedValue) {
		t.Errorf("Expected flag variable: %v, Actual: %v", expectedValue, *actualVar)
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
