package cli

import (
    "github.com/dailymuse/git-fit/util"
    "github.com/dailymuse/git-fit/transport"
    "github.com/cheggaaa/pb"
    "fmt"
)

type operationResponse struct {
    transport.ProgressMessage
    Path string
    Message string
}

func newOperationResponse(path string, progress transport.ProgressMessage) operationResponse {
    return operationResponse {
        ProgressMessage: progress,
        Path: path,
    }
}

func newErrorOperationResponse(path string, err error) operationResponse {
    return operationResponse {
        ProgressMessage: transport.NewErrorProgressMessage(err),
        Path: path,
    }
}

func pipeResponses(path string, sendFinal bool, fromChan chan transport.ProgressMessage, toChan chan operationResponse) error {
    for {
        progress := <- fromChan
        isFinal := progress.Err == transport.ErrProgressCompleted

        if !isFinal || (isFinal && sendFinal) {
            toChan <- newOperationResponse(path, progress)
        }

        if progress.Err != nil {
            return progress.Err
        }
    }
}

func handleResponse(ch chan operationResponse, fileCount int) {
    if fileCount == 0 {
        return
    }

    statuses := make(map[string]operationResponse)
    bar := pb.StartNew(fileCount * 100)

    //TODO: use floats for progress?
    for {
        res := <- ch
        statuses[res.Path] = res
        progress := 0
        doneCount := 0

        for _, status := range statuses {
            if status.Err != nil {
                progress += 100
                doneCount++
            } else {
                progress += int(status.Percent * 100)
            }
        }

        bar.Set(progress)

        if doneCount == fileCount {
            bar.Finish()

            for _, status := range statuses {
                if status.Err != transport.ErrProgressCompleted {
                    util.Error(fmt.Sprintf("%s: %s\n", status.Path, status.Err))
                } else {
                    util.Message(fmt.Sprintf("%s: Synced\n", status.Path))
                }
            }

            return
        }
    }
}
