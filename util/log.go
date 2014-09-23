package util

import (
    "os"
    "fmt"
)

func Error(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, format, args...)
}

func Fatal(format string, args ...interface{}) {
    Error(format, args...)
    os.Exit(1)
}

func Message(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}
