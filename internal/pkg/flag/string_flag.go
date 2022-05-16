package flag

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

func (f *StringFlag) Default(value string) {
	f.value = value
}

func (f StringFlag) String() string {
	return f.value
}

func (f StringFlag) Value() string {
	return f.value
}

func (f *StringFlag) Set(value string) error {
	f.value = value
	f.provided = true

	return nil
}

func (f StringFlag) Provided() bool {
	return f.provided
}

func (f *StringFlag) Type() string {
	return "string"
}
