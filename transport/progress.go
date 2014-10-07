package transport

import (
    "errors"
)

var ErrProgressCompleted error = errors.New("Progress done")

type ProgressMessage struct {
    Percent float64
    Err error
}

func NewProgressMessage(percent float64) ProgressMessage {
    return ProgressMessage {
        Percent: percent,
        Err: nil,
    }
}

func NewErrorProgressMessage(err error) ProgressMessage {
    return ProgressMessage {
        Percent: 1.0,
        Err: err,
    }
}

func NewFinishedProgressMessage() ProgressMessage {
    return ProgressMessage {
        Percent: 1.0,
        Err: ErrProgressCompleted,
    }
}

type ProgressWriter struct {
    ProgressChan chan ProgressMessage
    Written uint64
    TotalSize uint64
}

func NewProgressWriter(progressChan chan ProgressMessage, totalSize uint64) *ProgressWriter {
    writer := ProgressWriter {
        ProgressChan: progressChan,
        Written: 0,
        TotalSize: totalSize,
    }

    return &writer
}

func (self *ProgressWriter) Write(p []byte) (n int, err error) {
    self.Written += uint64(len(p))
    self.ProgressChan <- NewProgressMessage(float64(self.Written) / float64(self.TotalSize))
    return len(p), nil
}
