package transport

import (
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
