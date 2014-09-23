package cli

import (
    "os"
    "fmt"
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
    "errors"
)

func download(trans transport.Transport, file transport.RemotableFile, responseChan chan operationResponse) {
    state, err := getRemotableFileState(trans, file)

    if err != nil {
        responseChan <- newErrorOperationResponse(file, err)
    } else if state & REMOTABLE_FILE_STATE_NO_REMOTE != 0 {
        responseChan <- newErrorOperationResponse(file, errors.New("remote file does not exist"))
    } else if state & REMOTABLE_FILE_STATE_SAME != 0 {
        responseChan <- newOperationResponse(file, false)
    } else {
        if err = trans.Download(file); err != nil {
            responseChan <- newErrorOperationResponse(file, err)
        } else {
            responseChan <- newOperationResponse(file, true)
        }
    }
}

func Pull(schema config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema {
            if util.FileExists(path) {
                fmt.Fprintf(os.Stderr, "%s: Not overwriting because the file already exists. If you wish to overwrite the current contents, explicitly include the file path as an argument\n", path)
            } else {
                paths = append(paths, path)
            }
        }
    } else {
        for _, path := range paths {
            if _, ok := schema[path]; !ok {
                panic(fmt.Sprintf("%s: no entry in the config %s", path, config.FILE_PATH))
            }
        }
    }

    responseChan := make(chan operationResponse, len(paths))

    for _, path := range paths {
        remoteFile := transport.NewRemotableFile(schema[path].Commit, path)
        go download(trans, remoteFile, responseChan)
    }

    for i := 0; i < len(paths); i++ {
        res := <- responseChan

        if res.err != nil {
            fmt.Fprintf(os.Stderr, "%s: Could not download: %s\n", res.file.Path, res.err.Error())
        } else if !res.synced {
            fmt.Fprintf(os.Stderr, "%s: Already synced\n", res.file.Path)
        } else {
            fmt.Printf("%s: Downloaded\n", res.file.Path)
        }
    }
}
