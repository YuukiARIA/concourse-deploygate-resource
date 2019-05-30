package logger

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func Fatal(arg interface{}) {
	printMessage(arg)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	printMessage(aurora.Sprintf(aurora.Red(format), args...))
	os.Exit(1)
}

func Sayf(format string, args ...interface{}) {
	printMessage(fmt.Sprintf(format, args...))
}

func printMessage(message interface{}) {
	fmt.Fprintln(os.Stderr, message)
}
