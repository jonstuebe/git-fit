package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
)

func FileExists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func IsDirectory(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func IsFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func FileHash(p string) (string, error) {
	file, err := os.Open(p)
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
	err := os.MkdirAll(path.Dir(toPath), os.ModePerm)

	if err != nil {
		return err
	}

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
