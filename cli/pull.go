package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
    "errors"
)

func download(trans transport.Transport, path string, committedHash string, responseChan chan operationResponse) {
    if util.FileExists(path) {
        responseChan <- newOperationResponse(path, transport.NewProgressMessage(0, 0, errors.New("Skipped")), nil)
        return
    }

    blob := transport.NewBlob(committedHash)
    downloaded := false
    var progress transport.ProgressMessage

    if !util.IsFile(blob.Path()) {
        if progress = pipeProgress(path, trans.Download(blob), responseChan); progress.IsErrored() {
            responseChan <- newOperationResponse(path, progress, nil)
            return
        }

        downloaded = true
    }

    if err := util.CopyFile(blob.Path(), path); err != nil {
        if downloaded {
            responseChan <- newOperationResponse(path, transport.NewProgressMessage(progress.Total, progress.Total, err), nil)
        } else {
            responseChan <- newOperationResponse(path, transport.NewProgressMessage(0, 0, err), nil)
        }
    } else {
        if downloaded {
            responseChan <- newOperationResponse(path, progress, nil)
        } else {
            responseChan <- newOperationResponse(path, transport.NewProgressMessage(0, 0, transport.ErrProgressCompleted), nil)
        }
    }
}

func Pull(schema *config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema.Files {
            paths = append(paths, path)
        }
    } else {
        for _, path := range paths {
            if _, ok := schema.Files[path]; !ok {
                util.Fatal("%s: No entry in the config %s\n", path, config.FILE_PATH)
            }
        }
    }

    responseChan := make(chan operationResponse, len(paths))

    for _, path := range paths {
        go download(trans, path, schema.Files[path], responseChan)
    }

    handleResponse(responseChan, len(paths))
}
