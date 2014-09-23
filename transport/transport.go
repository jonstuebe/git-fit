package transport

import (
    "io"
    "os"
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "github.com/dailymuse/git-fit/util"
)

type RemotableFile struct {
    CommitHash string
    Path string
}

func NewRemotableFile(commitHash string, path string) RemotableFile {
    return RemotableFile {
        CommitHash: commitHash,
        Path: path,
    }
}

func (self RemotableFile) StandardName() string {
    return fmt.Sprintf("%s_%s", self.Path, self.CommitHash)
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

func (self RemotableFile) Hash() (string, error) {
    if !util.IsFile(self.Path) {
        return "", nil
    }

    file, err := self.GetFile()
    h := md5.New()

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
    LocalHash(file RemotableFile) (string, error)
    RemoteHash(file RemotableFile) (string, error)
}
