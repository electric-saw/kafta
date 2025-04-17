package flag

import (
	"strconv"

	"github.com/charmbracelet/huh"
)

// BoolFlag is a bool flag compatible with flags and pflags that keeps track of whether it had a value supplied or not.
type BoolFlag struct {
	// If Set has been invoked this value is true
	provided bool
	// The exact value provided on the flag
	value bool
}

func NewBool(defaultVal bool) BoolFlag {
	return BoolFlag{value: defaultVal}
}

func (f *BoolFlag) Default(value bool) {
	f.value = value
}

func (f *BoolFlag) Value() bool {
	return f.value
}

func (f *BoolFlag) String() string {
	return strconv.FormatBool(f.value)
}

func (f *BoolFlag) Set(value string) error {
	if bValue, err := strconv.ParseBool(value); err == nil {
		f.value = bValue
		f.provided = true
		return nil
	} else {
		return err
	}
}

func (f *BoolFlag) Provided() bool {
	return f.provided
}

func (f *BoolFlag) Type() string {
	return "bool"
}

func (f *BoolFlag) HuhWraps() huh.Accessor[bool] {
	return &huhWrapperBool{f: f}
}

type huhWrapperBool struct {
	f *BoolFlag
}

func (h *huhWrapperBool) Get() bool {
	return h.f.Value()
}

func (h *huhWrapperBool) Set(value bool) {
	h.f.provided = true
	h.f.value = value
}
