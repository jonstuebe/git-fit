package cli

import (
    "github.com/dailymuse/git-fit/transport"
)

type operationResponse struct {
    file transport.RemotableFile
    synced bool
    err error
}

func newOperationResponse(file transport.RemotableFile, synced bool) operationResponse {
    return operationResponse {
        file: file,
        synced: synced,
        err: nil,
    }
}

func newErrorOperationResponse(file transport.RemotableFile, err error) operationResponse {
    return operationResponse {
        file: file,
        synced: false,
        err: err,
    }
}
