package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
)

func upload(trans transport.Transport, file transport.RemotableFile, responseChan chan operationResponse) {
    actualHash, err := file.CalculateHash()

    if err != nil {
        responseChan <- newErrorOperationResponse(file, err)
    } else {
        actualFile := transport.NewRemotableFile(file.Path, actualHash)
        exists, err := trans.Exists(actualFile)

        if err != nil {
            responseChan <- newErrorOperationResponse(actualFile, err)
        } else if exists {
            responseChan <- newOperationResponse(actualFile, false)
        } else if err = trans.Upload(actualFile); err != nil {
            responseChan <- newErrorOperationResponse(actualFile, err)
        } else {
            responseChan <- newOperationResponse(actualFile, true)
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
        remoteFile := transport.NewRemotableFile(path, schema.Files[path])
        go upload(trans, remoteFile, responseChan)
    }

    for i := 0; i < len(paths); i++ {
        res := <- responseChan

        if res.err != nil {
            util.Error("%s: Could not upload: %s\n", res.file.Path, res.err.Error())
        } else if !res.synced {
            util.Error("%s: Already synced\n", res.file.Path)
        } else {
            util.Message("%s: Uploaded\n", res.file.Path)
            schema.Files[res.file.Path] = res.file.CommittedHash
        }
    }
}
