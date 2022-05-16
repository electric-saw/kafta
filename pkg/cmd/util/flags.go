package util

type StringFlag struct {
	provided bool
	value    string
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
