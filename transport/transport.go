package transport

import (
    "io"
    "os"
    "fmt"
    "crypto/sha1"
    "encoding/hex"
    "github.com/dailymuse/git-fit/util"
)

type RemotableFile struct {
    Path string
    CommittedHash string
}

func NewRemotableFile(path string, committedHash string) RemotableFile {
    return RemotableFile {
        Path: path,
        CommittedHash: committedHash,
    }
}

func (self RemotableFile) StandardName() string {
    return fmt.Sprintf("%s_%s", self.Path, self.CommittedHash)
}

func (self RemotableFile) GetFile() (*os.File, error) {
    return os.Open(self.Path)
}

func (self RemotableFile) WriteFile(source io.Reader) error {
    file, err := os.Create(self.Path)

    if err != nil {
        return err
    }

    defer file.Close()

    _, err = io.Copy(file, source)
    return err
}

func (self RemotableFile) CalculateHash() (string, error) {
    if !util.IsFile(self.Path) {
        return "", nil
    }

    file, err := self.GetFile()
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

type Transport interface {
    Upload(file RemotableFile) error
    Download(file RemotableFile) error
    Exists(file RemotableFile) (bool, error)
}
