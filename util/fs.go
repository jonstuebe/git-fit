package util

import (
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
