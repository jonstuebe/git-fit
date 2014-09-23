package cli

import (
    "os"
    "fmt"
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
    "errors"
)

func upload(trans transport.Transport, file transport.RemotableFile, responseChan chan operationResponse) {
    state, err := getRemotableFileState(trans, file)

    if err != nil {
        responseChan <- newErrorOperationResponse(file, err)
    } else if state & REMOTABLE_FILE_STATE_NO_LOCAL != 0 {
        responseChan <- newErrorOperationResponse(file, errors.New("does not exist"))
    } else if state & REMOTABLE_FILE_STATE_SAME != 0 {
        responseChan <- newOperationResponse(file, false)
    } else {
        if err = trans.Upload(file); err != nil {
            responseChan <- newErrorOperationResponse(file, err)
        } else {
            responseChan <- newOperationResponse(file, true)
        }
    }
}

func Push(schema config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema {
            paths = append(paths, path)
        }
    }

    latestCommit := util.LatestCommit()
    responseChan := make(chan operationResponse, len(args))

    for _, path := range args {
        remoteFile := transport.NewRemotableFile(latestCommit, path)
        go upload(trans, remoteFile, responseChan)
    }

    for i := 0; i < len(args); i++ {
        res := <- responseChan

        if res.err != nil {
            fmt.Fprintf(os.Stderr, "%s: Could not upload: %s\n", res.file.Path, res.err.Error())
        } else if !res.synced {
            fmt.Fprintf(os.Stderr, "%s: Already synced\n", res.file.Path)
        } else {
            fmt.Printf("%s: Uploaded\n", res.file.Path)

            schema[res.file.Path] = config.ConfigEntry {
                Commit: res.file.CommitHash,
            }
        }
    }
}
