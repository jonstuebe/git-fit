package transport

import (
    "io"
    "os"
    "fmt"
)

type Blob struct {
    Hash string
}

func NewBlob(hash string) Blob {
    return Blob {
        Hash: hash,
    }
}

func (self Blob) Path() string {
    return fmt.Sprintf(".git/fit/%s", self.Hash)
}

func (self Blob) Write(source io.Reader) error {
    file, err := os.Create(self.Path())

    if err != nil {
        return err
    }

    defer file.Close()

    _, err = io.Copy(file, source)
    return err
}

type Transport interface {
    Upload(blob Blob) error
    Download(blob Blob) error
    Exists(blob Blob) (bool, error)
}
