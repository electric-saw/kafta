package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		_, err := fmt.Fprint(os.Stderr, msg)
		if err != nil {
			fmt.Printf("Error writing to stderr: %v\n", err)
		}
	}
	os.Exit(code)
}

//nolint:gochecknoglobals // It is set to fatal by default, but can be overridden for testing.
var fatalErrHandler = fatal

const (
	DefaultErrorExitCode = 1
)

func UsageErrorf(cmd *cobra.Command, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}

	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = "error: " + msg
	}
	handleErr(msg, DefaultErrorExitCode)
}

func HelpErrorf(cmd *cobra.Command, format string, args ...any) error {
	CheckErr(cmd.Help())
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s", msg)
}

func HelpError(cmd *cobra.Command, args ...any) error {
	CheckErr(cmd.Help())
	msg := fmt.Sprint(args...)
	return fmt.Errorf("%s", msg)
}

// mapToMapPointer split string=string to a map[string]string.
func StringToMapPointer(s string) map[string]*string {
	m := make(map[string]*string)
	for v := range strings.SplitSeq(s, ",") {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			m[kv[0]] = &kv[1]
		}
	}
	return m
}
