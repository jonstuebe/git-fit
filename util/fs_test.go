package util

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	t.Parallel()

	if !FileExists("../README.md") {
		t.Error()
	}

	if FileExists("foo") {
		t.Error()
	}
}

func TestIsDirectory(t *testing.T) {
	t.Parallel()

	if !IsDirectory("../util") {
		t.Error()
	}

	if IsDirectory("foo") {
		t.Error()
	}
}

func TestFileHash(t *testing.T) {
	t.Parallel()

	hash, err := FileHash("../LICENSE")

	if err != nil {
		t.Error(err)
	}

	if hash != "96856a925efa12f5967f52455734b15fd2695e3a" {
		t.Error()
	}

	hash, err = FileHash("foo")

	if err == nil {
		t.Error()
	}

	if hash != "" {
		t.Error()
	}
}

func TestCopyFile(t *testing.T) {
	t.Parallel()

	err := CopyFile("from-path/does/not/exist", "to-path/does/not/exist")
	defer os.RemoveAll("to-path")

	if err == nil {
		t.Error("Expected an error")
	}

	if !IsDirectory("to-path") {
		t.Error("to-path does not exist")
	}

	err = CopyFile("../README.md", "temp")
	defer os.Remove("temp")

	if err != nil {
		t.Error(err)
	}
}
