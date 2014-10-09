package transport

import (
    "errors"
)

var ErrProgressCompleted error = errors.New("Progress done")

type ProgressMessage struct {
    Progress int
    Total int
    Err error
}

func NewProgressMessage(progress int, total int, err error) ProgressMessage {
    return ProgressMessage {
        Progress: progress,
        Total: total,
        Err: err,
    }
}

func (self ProgressMessage) IsErrored() bool {
    return self.Err != nil && self.Err != ErrProgressCompleted
}

func (self ProgressMessage) IsCompleted() bool {
    return self.Err == ErrProgressCompleted
}

type ProgressWriter struct {
    ProgressChan chan ProgressMessage
    Progress uint64
    Total uint64
}

func NewProgressWriter(progressChan chan ProgressMessage, total uint64) *ProgressWriter {
    writer := ProgressWriter {
        ProgressChan: progressChan,
        Progress: 0,
        Total: total,
    }

    return &writer
}

func (self *ProgressWriter) Write(p []byte) (n int, err error) {
    self.Progress += uint64(len(p))
    self.ProgressChan <- NewProgressMessage(int(self.Progress / 1024 / 1024), int(self.Total / 1024 / 1024), nil)
    return len(p), nil
}
