package transport

import (
    "errors"
    "testing"
)

func TestProgressMessageIsErrored(t *testing.T) {
    t.Parallel()

    progress := NewProgressMessage(0, 10, nil)

    if progress.IsErrored() {
        t.Error()
    }

    progress = NewProgressMessage(0, 10, ErrProgressCompleted)

    if progress.IsErrored() {
        t.Error()
    }

    progress = NewProgressMessage(0, 10, errors.New("Uhoh"))

    if !progress.IsErrored() {
        t.Error()
    }
}

func TestProgressMessageIsCompleted(t *testing.T) {
    t.Parallel()

    progress := NewProgressMessage(1, 1, nil)

    if progress.IsCompleted() {
        t.Error()
    }

    progress = NewProgressMessage(1, 1, ErrProgressCompleted)

    if !progress.IsCompleted() {
        t.Error()
    }

    progress = NewProgressMessage(1, 1, errors.New("Uhoh"))

    if progress.IsCompleted() {
        t.Error()
    }
}

func TestProgressWriterWrite(t *testing.T) {
    t.Parallel()

    ch := make(chan ProgressMessage, 1)
    writer := NewProgressWriter(ch, 10 * 1024 * 1024)
    n, err := writer.Write(make([]byte, 1 * 1024 * 1024))

    if err != nil {
        t.Error(err)
    }

    if n != 1 * 1024 * 1024 {
        t.Error()
    }

    progress := <- ch

    if progress.Progress != 1 {
        t.Error()
    }

    if progress.Total != 10 {
        t.Error()
    }

    if progress.Err != nil {
        t.Error(err)
    }
}
