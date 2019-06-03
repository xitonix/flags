package flags

import (
	"go.xitonix.io/flags/core"
	"reflect"
	"testing"

	"go.xitonix.io/flags/mocks"
)

func TestEnableAutoKeyGeneration(t *testing.T) {
	EnableAutoKeyGeneration()
	if !DefaultBucket.Options().AutoKeys {
		t.Error("The default bucket's auto key generation was expected to be ON")
	}
}

func TestSetKeyPrefix(t *testing.T) {
	prefix := "prefix"
	expected := "PREFIX"
	SetKeyPrefix(prefix)
	actual := DefaultBucket.Options().KeyPrefix
	if actual != expected {
		t.Errorf("The default bucket's key prefix was expected to be %s, but it was %s", expected, actual)
	}
}

func TestSetLogger(t *testing.T) {
	lg := &mocks.Logger{}
	SetLogger(lg)
	actual := DefaultBucket.Options().Logger
	if actual != lg {
		t.Error("The default bucket's logger has not been set as expected")
	}
}

func TestAppendSource(t *testing.T) {
	DefaultBucket = NewBucket()
	src := NewMemorySource()
	AppendSource(src)
	actual := DefaultBucket.sources[len(DefaultBucket.sources)-1]
	if actual != src {
		t.Error("The default bucket's source has not been appended as expected")
	}
}

func TestPrependSource(t *testing.T) {
	DefaultBucket = NewBucket()
	src := NewMemorySource()
	PrependSource(src)
	actual := DefaultBucket.sources[0]
	if actual != src {
		t.Error("The default bucket's source has not been prepended as expected")
	}
}

func TestAddSource(t *testing.T) {
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

			DefaultBucket = NewBucket()
			AddSource(tc.src, tc.index)

			if len(DefaultBucket.sources) != tc.expectedLength {
				t.Errorf("Expected Number of Sources: %d, Actual: %d", tc.expectedLength, len(DefaultBucket.sources))
				return
			}

			for i, expected := range tc.expected {
				actual := DefaultBucket.sources[i]
				if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
					t.Errorf("Expected Source at index %d: %T, Actual: %T", i, expected, actual)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	DefaultBucket = NewBucket()
	Options().Terminator = &mocks.Terminator{}
	Options().Logger = &mocks.Logger{}
	Options().HelpProvider.Writer = mocks.NewInMemoryWriter()
	String("long", "usage")
	Parse()
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
}

func TestOptions(t *testing.T) {
	tm := &mocks.Terminator{}
	Options().Terminator = tm
	if DefaultBucket.Options().Terminator != tm {
		t.Error("The default bucket's options has not been set as expected")
	}
}

func TestGlobalString(t *testing.T) {
	DefaultBucket = NewBucket()
	String("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringFlag{}, f)
	}
}

func TestGlobalStringP(t *testing.T) {
	DefaultBucket = NewBucket()
	StringP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringFlag{}, f)
	}
}

func TestGlobalInt(t *testing.T) {
	DefaultBucket = NewBucket()
	Int("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntFlag{}, f)
	}
}

func TestGlobalIntP(t *testing.T) {
	DefaultBucket = NewBucket()
	IntP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntFlag{}, f)
	}
}

func TestGlobalInt64(t *testing.T) {
	DefaultBucket = NewBucket()
	Int64("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int64Flag{}, f)
	}
}

func TestGlobalInt64P(t *testing.T) {
	DefaultBucket = NewBucket()
	Int64P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int64Flag{}, f)
	}
}
