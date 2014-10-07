package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
    "errors"
)

func upload(trans transport.Transport, path string, responseChan chan operationResponse) {
    hash, err := util.FileHash(path)

    if err != nil {
        responseChan <- newErrorOperationResponse(path, err)
    } else if hash == "" {
        responseChan <- newErrorOperationResponse(path, errors.New("File does not exist"))
    } else {
        blob := transport.NewBlob(hash)

        if !util.IsFile(blob.Path()) {
            if err = util.CopyFile(path, blob.Path()); err != nil {
                responseChan <- newErrorOperationResponse(path, err)
                return
            }
        }

        exists, err := trans.Exists(blob)

        if err != nil {
            responseChan <- newErrorOperationResponse(path, err)
        } else if exists {
            responseChan <- newOperationResponse(path, transport.NewFinishedProgressMessage())
        } else {
            pipeResponses(path, true, trans.Upload(blob), responseChan)
        }
    }
}

func Push(schema *config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema.Files {
            paths = append(paths, path)
        }
    }

    responseChan := make(chan operationResponse, len(paths))

    for _, path := range paths {
        go upload(trans, path, responseChan)
    }

    handleResponse(responseChan, len(paths))
}
