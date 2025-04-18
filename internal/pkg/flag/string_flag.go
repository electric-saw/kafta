package flag

import "github.com/charmbracelet/huh"

// StringFlag is a string flag compatible with flags and pflags that keeps track of whether it had a value supplied or not.
type StringFlag struct {
	// If Set has been invoked this value is true
	provided bool
	// The exact value provided on the flag
	value string
}

func NewStringFlag(defaultVal string) StringFlag {
	return StringFlag{value: defaultVal}
}

func (f *StringFlag) String() string {
	return f.value
}

func (f *StringFlag) Value() string {
	return f.value
}

func (f *StringFlag) Set(value string) error {
	f.value = value
	f.provided = true

	return nil
}

func (f *StringFlag) Provided() bool {
	return f.provided
}

func (f *StringFlag) Type() string {
	return "string"
}

func (f *StringFlag) HuhWraps() huh.Accessor[string] {
	return &huhWrapperString{f: f}
}

type huhWrapperString struct {
	f *StringFlag
}

func (h *huhWrapperString) Get() string {
	return h.f.Value()
}

func (h *huhWrapperString) Set(value string) {
	_ = h.f.Set(value) //nolint:errcheck // this err is always nil
}
