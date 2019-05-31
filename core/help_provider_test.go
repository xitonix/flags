package core

import (
	"testing"

	"go.xitonix.io/flags/mocks"
)

func TestNewHelpProvider(t *testing.T) {
	w := mocks.NewInMemoryWriter()
	tf := &TabbedHelpFormatter{}
	h := NewHelpProvider(w, tf)
	if h.Formatter != tf {
		t.Error("Tab formatter has not been set properly")
	}
	if h.Writer != w {
		t.Error("Help writer has not been set properly")
	}
}
