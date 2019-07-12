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
	if _, ok := f.(*core.StringFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.StringFlag{}, f)
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
	if _, ok := f.(*core.StringMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.StringMapFlag{}, f)
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
	if _, ok := f.(*core.StringSliceMapFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.StringSliceMapFlag{}, f)
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
	if _, ok := f.(*core.IntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.IntFlag{}, f)
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
	if _, ok := f.(*core.Int64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Int64Flag{}, f)
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
	if _, ok := f.(*core.Int32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Int32Flag{}, f)
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
	if _, ok := f.(*core.Int16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Int16Flag{}, f)
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
	if _, ok := f.(*core.Int8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Int8Flag{}, f)
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
	if _, ok := f.(*core.UIntFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UIntFlag{}, f)
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
	if _, ok := f.(*core.UInt64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UInt64Flag{}, f)
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
	if _, ok := f.(*core.UInt32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UInt32Flag{}, f)
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
	if _, ok := f.(*core.UInt16Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UInt16Flag{}, f)
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
	if _, ok := f.(*core.UInt8Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UInt8Flag{}, f)
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
	if _, ok := f.(*core.ByteFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.ByteFlag{}, f)
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
	if _, ok := f.(*core.Float64Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Float64Flag{}, f)
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
	if _, ok := f.(*core.Float32Flag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Float32Flag{}, f)
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
	if _, ok := f.(*core.CounterFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.CounterFlag{}, f)
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
	if _, ok := f.(*core.CounterFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.CounterFlag{}, f)
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
	if _, ok := f.(*core.DurationFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.DurationFlag{}, f)
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
	if _, ok := f.(*core.DurationSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.DurationSliceFlag{}, f)
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
	if _, ok := f.(*core.BoolFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.BoolFlag{}, f)
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
	if _, ok := f.(*core.BoolSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.BoolSliceFlag{}, f)
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
	if _, ok := f.(*core.TimeFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.TimeFlag{}, f)
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
	if _, ok := f.(*core.StringSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.StringSliceFlag{}, f)
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
	if _, ok := f.(*core.IntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.IntSliceFlag{}, f)
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
	if _, ok := f.(*core.UIntSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.UIntSliceFlag{}, f)
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
	if _, ok := f.(*core.Float64SliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.Float64SliceFlag{}, f)
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
	if _, ok := f.(*core.IPAddressFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.IPAddressFlag{}, f)
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
	if _, ok := f.(*core.IPAddressSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.IPAddressSliceFlag{}, f)
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
	if _, ok := f.(*core.CIDRFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.CIDRFlag{}, f)
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
	if _, ok := f.(*core.CIDRSliceFlag); !ok {
		t.Errorf("Expected %T, but received %T", &core.CIDRSliceFlag{}, f)
	}
}
