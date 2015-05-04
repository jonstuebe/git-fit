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

	err := CopyFile("../README.md", "temp")

	if err != nil {
		t.Error(err)
	}

	err = os.Remove("temp")

	if err != nil {
		t.Fatal(err)
	}

	err = CopyFile("temp", "temp2")

	if err == nil {
		os.Remove("temp2")
		t.Error(err)
	}
}
