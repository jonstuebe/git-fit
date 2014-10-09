package util

import (
    "testing"
)

func TestStdout(t *testing.T) {
    t.Parallel()

    results := Stdout("echo", "hi")

    if results != "hi" {
        t.Fail()
    }
}

// For the next two functions, just make sure they don't panic, as they could
// have different outputs depending on the context
func TestGitDir(t *testing.T) {
    t.Parallel()
    GitDir()
}

func TestGitConfig(t *testing.T) {
    t.Parallel()
    GitConfig("foo")
}
