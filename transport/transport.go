package transport

type Transport interface {
    Upload(blob Blob) (chan ProgressMessage)
    Download(blob Blob) (chan ProgressMessage)
    Exists(blob Blob) (bool, error)
}
