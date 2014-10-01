package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
    "github.com/dailymuse/git-fit/util"
)

func download(trans transport.Transport, path string, committedHash string, responseChan chan operationResponse) {
    blob := transport.NewBlob(committedHash)
    synced := false

    if !util.IsFile(blob.Path()) {
        if err := trans.Download(blob); err != nil {
            responseChan <- newErrorOperationResponse(path, err)
            return
        }

        synced = true
    }

    if err := util.CopyFile(blob.Path(), path); err != nil {
        responseChan <- newErrorOperationResponse(path, err)   
    } else {
        responseChan <- newOperationResponse(path, synced)
    }
}

func Pull(schema *config.Config, trans transport.Transport, args []string) {
    paths := args

    if len(paths) == 0 {
        paths = make([]string, 0)

        for path := range schema.Files {
            if util.FileExists(path) {
                util.Error("%s: Not overwriting because the file already exists. If you wish to overwrite the current contents, explicitly include the file path as an argument.\n", path)
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

    for i := 0; i < len(paths); i++ {
        res := <- responseChan

        if res.err != nil {
            util.Error("%s: Could not download: %s\n", res.path, res.err.Error())
        } else {
            synced := res.response.(bool)

            if !synced {
                util.Error("%s: Already synced\n", res.path)
            } else {
                util.Message("%s: Downloaded\n", res.path)
            }
        }
    }
}
