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

func pipeResponses(path string, sendFinal bool, fromChan chan transport.ProgressMessage, toChan chan operationResponse) transport.ProgressMessage {
    for {
        progress := <- fromChan

        if !progress.IsCompleted() || (progress.IsCompleted() && sendFinal) {
            toChan <- newOperationResponse(path, progress)
        }

        if progress.Err != nil {
            return progress
        }
    }
}

func handleResponse(ch chan operationResponse, fileCount int) {
    if fileCount == 0 {
        return
    }

    statuses := make(map[string]operationResponse)
    var bar *pb.ProgressBar

    for {
        res := <- ch
        statuses[res.Path] = res

        if len(statuses) == fileCount {
            total := 0
            progress := 0
            doneCount := 0

            for _, status := range statuses {
                progress += status.Progress
                total += status.Total

                if status.Err != nil {
                    doneCount++
                }
            }

            if bar == nil && progress > 0 {
                bar = pb.StartNew(total)
            }

            if bar != nil {
                bar.Set(progress)
            }

            if doneCount == fileCount {
                if bar != nil {
                    bar.Finish()
                }

                for _, status := range statuses {
                    if status.IsErrored() {
                        util.Error(fmt.Sprintf("%s: %s\n", status.Path, status.Err))
                    } else {
                        util.Message(fmt.Sprintf("%s: Synced\n", status.Path))
                    }
                }

                return
            }
        }
    }
}
