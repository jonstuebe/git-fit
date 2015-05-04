package cli

import (
	"github.com/cheggaaa/pb"
	"github.com/dailymuse/git-fit/transport"
	"github.com/dailymuse/git-fit/util"
)

type operationResponse struct {
	transport.ProgressMessage
	Path    string
	Payload interface{}
}

func newOperationResponse(path string, progress transport.ProgressMessage, payload interface{}) operationResponse {
	return operationResponse{
		ProgressMessage: progress,
		Path:            path,
		Payload:         payload,
	}
}

func pipeProgress(path string, fromChan chan transport.ProgressMessage, toChan chan operationResponse) transport.ProgressMessage {
	for {
		progress := <-fromChan

		if !progress.IsCompleted() {
			toChan <- newOperationResponse(path, progress, nil)
		}

		if progress.Err != nil {
			return progress
		}
	}
}

func handleResponse(ch chan operationResponse, fileCount int) []operationResponse {
	if fileCount == 0 {
		return make([]operationResponse, 0)
	}

	statuses := make(map[string]operationResponse)
	var bar *pb.ProgressBar

	for {
		res := <-ch
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

				successful := make([]operationResponse, 0)

				for _, status := range statuses {
					if status.IsErrored() {
						util.Error("%s: %s\n", status.Path, status.Err)
					} else {
						util.Message("%s: Synced\n", status.Path)
						successful = append(successful, status)
					}
				}

				return successful
			}
		}
	}
}
