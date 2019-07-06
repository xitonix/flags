package flags

import (
	"reflect"
	"testing"

	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/mocks"
)

func TestEnableAutoKeyGeneration(t *testing.T) {
	EnableAutoKeyGeneration()
	if !DefaultBucket.opts.AutoKeys {
		t.Error("The default bucket's auto key generation was expected to be ON")
	}
}

func TestSetKeyPrefix(t *testing.T) {
	prefix := "prefix"
	expected := "PREFIX"
	SetKeyPrefix(prefix)
	actual := DefaultBucket.opts.KeyPrefix
	if actual != expected {
		t.Errorf("The default bucket's key prefix was expected to be %s, but it was %s", expected, actual)
	}
}

func TestSetLogger(t *testing.T) {
	lg := &mocks.Logger{}
	SetLogger(lg)
	actual := DefaultBucket.opts.Logger
	if actual != lg {
		t.Error("The default bucket's logger has not been set as expected")
	}
}

func TestSetSortOrder(t *testing.T) {
	expected := by.KeyAscending
	SetSortOrder(expected)
	actual := DefaultBucket.opts.Comparer
	if actual != expected {
		t.Errorf("The default bucket's sort order was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetHelpFormatter(t *testing.T) {
	expected := &core.TabbedHelpFormatter{}
	SetHelpFormatter(expected)
	actual := DefaultBucket.opts.HelpFormatter
	if actual != expected {
		t.Errorf("The default bucket's help formatter was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetHelpWriter(t *testing.T) {
	expected := mocks.NewInMemoryWriter()
	SetHelpWriter(expected)
	actual := DefaultBucket.opts.HelpWriter
	if actual != expected {
		t.Errorf("The default bucket's help writer was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetTerminator(t *testing.T) {
	expected := &mocks.Terminator{}
	SetTerminator(expected)
	actual := DefaultBucket.opts.Terminator
	if actual != expected {
		t.Errorf("The default bucket's terminator was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetDeprecationMark(t *testing.T) {
	expected := "Deprecation Mark"
	SetDeprecationMark(expected)
	actual := DefaultBucket.opts.DeprecationMark
	if actual != expected {
		t.Errorf("The default bucket's deprecation mark was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetRequiredMark(t *testing.T) {
	expected := "**"
	SetRequiredFlagMark(expected)
	actual := DefaultBucket.opts.RequiredFlagMark
	if actual != expected {
		t.Errorf("The default bucket's Required mark was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetDefaultValueFormatString(t *testing.T) {
	expected := "Format FullString"
	SetDefaultValueFormatString(expected)
	actual := DefaultBucket.opts.DefaultValueFormatString
	if actual != expected {
		t.Errorf("The default bucket's default value format string was expected to be %T, but it was %T", expected, actual)
	}
}

func TestSetPreSetCallback(t *testing.T) {
	counter := 0
	expected := func(f core.Flag, value string) error {
		if value != "value" {
			t.Errorf("Expected value: 'value', Actual: %s", value)
		}
		counter++
		return nil
	}
	SetPreSetCallback(expected)
	if DefaultBucket.opts.PreSetCallback == nil {
		t.Error("The default bucket's pre-set callback was nil")
		return
	}
	_ = DefaultBucket.opts.PreSetCallback(mocks.NewFlag("long", "s"), "value")
	if counter != 1 {
		t.Error("The default bucket's pre-set callback has not been called")
	}
}

func TestSetPostSetCallback(t *testing.T) {
	counter := 0
	expected := func(f core.Flag, value string) error {
		if value != "value" {
			t.Errorf("Expected value: 'value', Actual: %s", value)
		}
		counter++
		return nil
	}
	SetPostSetCallback(expected)
	if DefaultBucket.opts.PostSetCallback == nil {
		t.Error("The default bucket's post-set callback was nil")
		return
	}
	_ = DefaultBucket.opts.PostSetCallback(mocks.NewFlag("long", "s"), "value")
	if counter != 1 {
		t.Error("The default bucket's post-set callback has not been called")
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

func TestAddFlag(t *testing.T) {
	f := mocks.NewFlag("long", "short")
	Add(f)
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	af := DefaultBucket.Flags()[0]
	if _, ok := af.(*mocks.Flag); !ok {
		t.Errorf("Expected %T, but received %T", &mocks.Flag{}, af)
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
	DefaultBucket.Options().Terminator = &mocks.Terminator{}
	DefaultBucket.Options().Logger = &mocks.Logger{}
	DefaultBucket.Options().HelpWriter = mocks.NewInMemoryWriter()
	String("long", "usage")
	Parse()
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
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

func TestGlobalStringMap(t *testing.T) {
	DefaultBucket = NewBucket()
	StringMap("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringMapFlag{}, f)
	}
}

func TestGlobalStringMapP(t *testing.T) {
	DefaultBucket = NewBucket()
	StringMapP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringMapFlag{}, f)
	}
}

func TestGlobalStringSliceMap(t *testing.T) {
	DefaultBucket = NewBucket()
	StringSliceMap("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringSliceMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceMapFlag{}, f)
	}
}

func TestGlobalStringSliceMapP(t *testing.T) {
	DefaultBucket = NewBucket()
	StringSliceMapP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringSliceMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceMapFlag{}, f)
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

func TestGlobalInt32(t *testing.T) {
	DefaultBucket = NewBucket()
	Int32("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int32Flag{}, f)
	}
}

func TestGlobalInt32P(t *testing.T) {
	DefaultBucket = NewBucket()
	Int32P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int32Flag{}, f)
	}
}

func TestGlobalInt16(t *testing.T) {
	DefaultBucket = NewBucket()
	Int16("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int16Flag{}, f)
	}
}

func TestGlobalInt16P(t *testing.T) {
	DefaultBucket = NewBucket()
	Int16P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int16Flag{}, f)
	}
}

func TestGlobalInt8(t *testing.T) {
	DefaultBucket = NewBucket()
	Int8("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int8Flag{}, f)
	}
}

func TestGlobalInt8P(t *testing.T) {
	DefaultBucket = NewBucket()
	Int8P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Int8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Int8Flag{}, f)
	}
}

func TestGlobalUInt(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UIntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntFlag{}, f)
	}
}

func TestGlobalUIntP(t *testing.T) {
	DefaultBucket = NewBucket()
	UIntP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UIntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntFlag{}, f)
	}
}

func TestGlobalUInt64(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt64("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt64Flag{}, f)
	}
}

func TestGlobalUInt64P(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt64P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt64Flag{}, f)
	}
}

func TestGlobalUInt32(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt32("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt32Flag{}, f)
	}
}

func TestGlobalUInt32P(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt32P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt32Flag{}, f)
	}
}

func TestGlobalUInt16(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt16("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt16Flag{}, f)
	}
}

func TestGlobalUInt16P(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt16P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt16Flag{}, f)
	}
}

func TestGlobalUInt8(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt8("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt8Flag{}, f)
	}
}

func TestGlobalUInt8P(t *testing.T) {
	DefaultBucket = NewBucket()
	UInt8P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UInt8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &UInt8Flag{}, f)
	}
}

func TestGlobalByte(t *testing.T) {
	DefaultBucket = NewBucket()
	Byte("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*ByteFlag); !ok {
		t.Errorf("Expected %T, but received %T", &ByteFlag{}, f)
	}
}

func TestGlobalFloat64(t *testing.T) {
	DefaultBucket = NewBucket()
	Float64("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64Flag{}, f)
	}
}

func TestGlobalFloat64P(t *testing.T) {
	DefaultBucket = NewBucket()
	Float64P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64Flag{}, f)
	}
}

func TestGlobalFloat32(t *testing.T) {
	DefaultBucket = NewBucket()
	Float32("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float32Flag{}, f)
	}
}

func TestGlobalFloat32P(t *testing.T) {
	DefaultBucket = NewBucket()
	Float32P("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &Float32Flag{}, f)
	}
}

func TestGlobalCounter(t *testing.T) {
	DefaultBucket = NewBucket()
	Counter("long", "usage").WithShort("c")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*CounterFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CounterFlag{}, f)
	}
}

func TestGlobalVerbosity(t *testing.T) {
	DefaultBucket = NewBucket()
	Verbosity("usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
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

func TestGlobalDuration(t *testing.T) {
	DefaultBucket = NewBucket()
	Duration("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*DurationFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationFlag{}, f)
	}
}

func TestGlobalDurationSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	DurationSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*DurationSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationSliceFlag{}, f)
	}
}

func TestGlobalDurationSliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	DurationSliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*DurationSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &DurationSliceFlag{}, f)
	}
}

func TestGlobalBool(t *testing.T) {
	DefaultBucket = NewBucket()
	Bool("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*BoolFlag); !ok {
		t.Errorf("Expected %T, but received %T", &BoolFlag{}, f)
	}
}

func TestGlobalBoolSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	BoolSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*BoolSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &BoolSliceFlag{}, f)
	}
}

func TestGlobalTime(t *testing.T) {
	DefaultBucket = NewBucket()
	Time("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*TimeFlag); !ok {
		t.Errorf("Expected %T, but received %T", &TimeFlag{}, f)
	}
}

func TestGlobalTimeP(t *testing.T) {
	DefaultBucket = NewBucket()
	TimeP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*TimeFlag); !ok {
		t.Errorf("Expected %T, but received %T", &TimeFlag{}, f)
	}
}

func TestGlobalStringSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	StringSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceFlag{}, f)
	}
}

func TestGlobalStringSliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	StringSliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*StringSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &StringSliceFlag{}, f)
	}
}

func TestGlobalIntSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	IntSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntSliceFlag{}, f)
	}
}

func TestGlobalIntSliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	IntSliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IntSliceFlag{}, f)
	}
}

func TestGlobalUIntSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	UIntSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UIntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntSliceFlag{}, f)
	}
}

func TestGlobalUIntSliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	UIntSliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*UIntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &UIntSliceFlag{}, f)
	}
}

func TestGlobalFloat64Slice(t *testing.T) {
	DefaultBucket = NewBucket()
	Float64Slice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float64SliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64SliceFlag{}, f)
	}
}

func TestGlobalFloat64SliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	Float64SliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*Float64SliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &Float64SliceFlag{}, f)
	}
}

func TestGlobalIPAddress(t *testing.T) {
	DefaultBucket = NewBucket()
	IPAddress("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IPAddressFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressFlag{}, f)
	}
}

func TestGlobalIPAddressP(t *testing.T) {
	DefaultBucket = NewBucket()
	IPAddressP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IPAddressFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressFlag{}, f)
	}
}

func TestGlobalIPAddressSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	IPAddressSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IPAddressSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressSliceFlag{}, f)
	}
}

func TestGlobalIPAddressSliceP(t *testing.T) {
	DefaultBucket = NewBucket()
	IPAddressSliceP("long", "s", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*IPAddressSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &IPAddressSliceFlag{}, f)
	}
}

func TestGlobalCIDR(t *testing.T) {
	DefaultBucket = NewBucket()
	CIDR("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*CIDRFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CIDRFlag{}, f)
	}
}

func TestGlobalCIDRSlice(t *testing.T) {
	DefaultBucket = NewBucket()
	CIDRSlice("long", "usage")
	actual := len(DefaultBucket.Flags())
	if actual != 1 {
		t.Errorf("Expected to get 1 parsed flag, but received %d", actual)
	}
	f := DefaultBucket.Flags()[0]
	if _, ok := f.(*CIDRSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &CIDRSliceFlag{}, f)
	}
}
