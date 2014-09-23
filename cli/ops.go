package cli

import (
    "github.com/dailymuse/git-fit/transport"
)

const (
    REMOTABLE_FILE_STATE_NONE = 0
    REMOTABLE_FILE_STATE_SAME = 1
    REMOTABLE_FILE_STATE_NO_LOCAL = 2
    REMOTABLE_FILE_STATE_NO_REMOTE = 4
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

func getRemotableFileState(trans transport.Transport, file transport.RemotableFile) (int, error) {
    localHash, err := trans.LocalHash(file)

    if err != nil {
        return REMOTABLE_FILE_STATE_NONE, err
    }

    remoteHash, err := trans.RemoteHash(file)

    if err != nil {
        return REMOTABLE_FILE_STATE_NONE, err
    }

    state := REMOTABLE_FILE_STATE_NONE

    if localHash == "" {
        state |= REMOTABLE_FILE_STATE_NO_LOCAL
    }

    if remoteHash == "" {
        state |= REMOTABLE_FILE_STATE_NO_REMOTE
    }

    if localHash == remoteHash {
        state |= REMOTABLE_FILE_STATE_SAME
    }

    return state, nil
}
