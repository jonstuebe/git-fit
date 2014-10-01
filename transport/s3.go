package transport

import (
    "os"
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

func (self S3Transport) Download(blob Blob) error {
    reader, err := self.bucket.GetReader(blob.Hash)

    if err != nil {
        return err
    }

    defer reader.Close()
    return blob.Write(reader)
}

func (self S3Transport) Upload(blob Blob) error {
    contents, err := os.Open(blob.Path())

    if err != nil {
        return err
    }

    defer contents.Close()

    multi, err := self.bucket.InitMulti(blob.Hash, "application/octet-stream", s3.Private)

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

func (self S3Transport) Exists(blob Blob) (bool, error) {
    _, err := self.bucket.GetKey(blob.Hash)

    // TODO: should be a better way of checking this
    if err != nil && err.Error() == "404 Not Found" {
        return false, nil
    } else {
        return err == nil, err
    }
}
