package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
)

func download(trans transport.Transport, path string, committedHash string, responseChan chan operationResponse) {
    blob := transport.NewBlob(committedHash)

    if !util.IsFile(blob.Path()) {
        err := pipeResponses(path, false, trans.Download(blob), responseChan); if err != transport.ErrProgressCompleted {
            return
        }
    }

    if err := util.CopyFile(blob.Path(), path); err != nil {
        responseChan <- newErrorOperationResponse(path, err)   
    } else {
        responseChan <- newOperationResponse(path, transport.NewFinishedProgressMessage())
    }
}

func Pull(schema *config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema.Files {
            if util.FileExists(path) {
                util.Error("%s: Skipped\n", path)
            } else {
                paths = append(paths, path)
            }
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
