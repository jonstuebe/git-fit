package transport

import (
    "github.com/mitchellh/goamz/s3"
)

const MULTIPART_CHUNK_SIZE = 8 * 1024 * 1024

type S3Transport struct {
    bucket *s3.Bucket
}

func NewS3Transport(bucket *s3.Bucket) S3Transport {
    return S3Transport {
        bucket: bucket,
    }
}

func (self S3Transport) Download(file RemotableFile) error {
    reader, err := self.bucket.GetReader(file.CommittedHash)

    if err != nil {
        return err
    }

    defer reader.Close()
    return file.WriteFile(reader)
}

func (self S3Transport) Upload(file RemotableFile) error {
    contents, err := file.GetFile()

    if err != nil {
        return err
    }

    defer contents.Close()

    multi, err := self.bucket.InitMulti(file.CommittedHash, "application/octet-stream", s3.Private)

    if err != nil {
        return err
    }

    parts, err := multi.PutAll(contents, MULTIPART_CHUNK_SIZE)

    if err != nil {
        multi.Abort()
        return err
    }

    return multi.Complete(parts)
}

func (self S3Transport) Exists(file RemotableFile) (bool, error) {
    _, err := self.bucket.GetKey(file.CommittedHash)

    // TODO: should be a better way of checking this
    if err != nil && err.Error() == "404 Not Found" {
        return false, nil
    } else {
        return err == nil, err
    }
}
