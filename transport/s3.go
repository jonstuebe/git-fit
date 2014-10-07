package transport

import (
    "io"
    "os"
    "bytes"
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

func (self S3Transport) downloadChunks(progress chan ProgressMessage, totalSize int64, file *os.File, reader io.ReadCloser) {
    defer file.Close()
    defer reader.Close()

    progressWriter := NewProgressWriter(progress, uint64(totalSize))
    writer := io.MultiWriter(progressWriter, file)
    _, err := io.Copy(writer, reader)

    if err != nil {
        progress <- NewErrorProgressMessage(err)
    } else {
        progress <- NewFinishedProgressMessage()
    }
}

func (self S3Transport) Download(blob Blob) (chan ProgressMessage) {
    progress := make(chan ProgressMessage, 10)
    key, err := self.bucket.GetKey(blob.Hash)

    if err != nil {
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    file, err := os.Create(blob.Path())

    if err != nil {
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    reader, err := self.bucket.GetReader(blob.Hash)

    if err != nil {
        file.Close()
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    go self.downloadChunks(progress, key.Size, file, reader)
    return progress
}

func (self S3Transport) uploadChunks(progress chan ProgressMessage, contents *os.File, info os.FileInfo, multi *s3.Multi) {
    defer contents.Close()

    chunk := make([]byte, MULTIPART_CHUNK_SIZE)
    totalChunks := info.Size() / MULTIPART_CHUNK_SIZE
    chunkNum := 0
    parts := make([]s3.Part, totalChunks)
    progress <- NewProgressMessage(0.0)

    for {
        n, bufferErr := io.ReadFull(contents, chunk)

        if bufferErr != nil && bufferErr != io.ErrUnexpectedEOF {
            progress <- NewErrorProgressMessage(bufferErr)
            multi.Abort()
            return
        }

        if n > 0 {
            reader := bytes.NewReader(chunk[:n])
            part, err := multi.PutPart(chunkNum, reader)

            if err != nil {
                progress <- NewErrorProgressMessage(err)
                multi.Abort()
                return
            }

            parts[chunkNum] = part
            chunkNum++
            progress <- NewProgressMessage(float64(chunkNum) / float64(totalChunks))
        }

        if bufferErr == io.ErrUnexpectedEOF {
            break
        }
    }

    if err := multi.Complete(parts); err != nil {
        progress <- NewErrorProgressMessage(err)
    } else {
        progress <- NewFinishedProgressMessage()
    }
}

func (self S3Transport) Upload(blob Blob) (chan ProgressMessage) {
    progress := make(chan ProgressMessage, 10)
    contents, err := os.Open(blob.Path())

    if err != nil {
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    info, err := contents.Stat()

    if err != nil {
        contents.Close()
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    multi, err := self.bucket.InitMulti(blob.Hash, "application/octet-stream", s3.Private)

    if err != nil {
        contents.Close()
        progress <- NewErrorProgressMessage(err)
        return progress
    }

    go self.uploadChunks(progress, contents, info, multi)
    return progress
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
