package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func FileHash(path string) (string, error) {
	file, err := os.Open(path)
	h := sha1.New()

	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err = io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func CopyFile(fromPath string, toPath string) error {
	in, err := os.Open(fromPath)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()

	if err != nil {
		return err
	}

	return cerr
}
