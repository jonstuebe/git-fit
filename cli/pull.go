package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
)

func download(trans transport.Transport, file transport.RemotableFile, responseChan chan operationResponse) {
    actualHash, err := file.CalculateHash()

    if err != nil {
        responseChan <- newErrorOperationResponse(file, err)
    } else if actualHash == file.CommittedHash {
        responseChan <- newOperationResponse(file, false)
    } else {
        if err = trans.Download(file); err != nil {
            responseChan <- newErrorOperationResponse(file, err)
        } else {
            responseChan <- newOperationResponse(file, true)
        }
    }
}

func Pull(schema *config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema.Files {
            if util.FileExists(path) {
                util.Error("%s: Not overwriting because the file already exists. If you wish to overwrite the current contents, explicitly include the file path as an argument\n", path)
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
        remoteFile := transport.NewRemotableFile(path, schema.Files[path])
        go download(trans, remoteFile, responseChan)
    }

    for i := 0; i < len(paths); i++ {
        res := <- responseChan

        if res.err != nil {
            util.Error("%s: Could not download: %s\n", res.file.Path, res.err.Error())
        } else if !res.synced {
            util.Error("%s: Already synced\n", res.file.Path)
        } else {
            util.Message("%s: Downloaded\n", res.file.Path)
        }
    }
}
