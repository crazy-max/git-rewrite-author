package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

// CheckIfError should be used to natively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}
	Error(err.Error())
}

// Error should be used to natively panics
func Error(format string, args ...interface{}) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", fmt.Sprintf(format, args...)))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// ParseAddress parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
func ParseAddress(address string) (name string, email string, err error) {
	re := regexp.MustCompile("(.*) <(.*)>")
	match := re.FindStringSubmatch(address)
	if match == nil || len(match) != 3 {
		return "", "", errors.New(fmt.Sprintf("Cannot parse %s", address))
	}
	return match[1], match[2], nil
}
