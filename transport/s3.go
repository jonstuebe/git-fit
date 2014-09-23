package transport

import (
    "github.com/mitchellh/goamz/s3"
    "io"
    "fmt"
    "crypto/md5"
    "encoding/hex"
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
    reader, err := self.bucket.GetReader(file.StandardName())

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

    multi, err := self.bucket.InitMulti(file.StandardName(), "application/octet-stream", s3.Private)

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

func (self S3Transport) LocalHash(file RemotableFile) (string, error) {
    contents, err := file.GetFile()

    if err != nil {
        //return "", err
        panic(err)
    }

    overallHasher := md5.New()
    chunkCount := 0

    defer contents.Close()

    for {
        chunk := make([]byte, MULTIPART_CHUNK_SIZE)
        n, err := io.ReadFull(contents, chunk)

        if err != nil && err != io.ErrUnexpectedEOF {
            //return "", err
            panic(err)
        }

        if n > 0 {
            chunkHasher := md5.New()
            chunkHasher.Write(chunk[0:n])
            overallHasher.Write(chunkHasher.Sum(nil))
            chunkCount++
        }

        if err == io.ErrUnexpectedEOF {
            break
        }
    }

    overallHash := hex.EncodeToString(overallHasher.Sum(nil))
    return fmt.Sprintf("%s-%d", overallHash, chunkCount), nil
}

func (self S3Transport) RemoteHash(file RemotableFile) (string, error) {
    key, err := self.bucket.GetKey(file.StandardName())

    if err != nil {
        // TODO: should be a better way of checking for a 404
        if err.Error() == "404 Not Found" {
            return "", nil
        } else {
            return "", err
        }
    }

    return key.ETag[1:len(key.ETag)-1], nil
}
