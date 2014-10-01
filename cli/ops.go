package cli

type operationResponse struct {
    path string
    response interface{}
    err error
}

func newOperationResponse(path string, response interface{}) operationResponse {
    return operationResponse {
        path: path,
        response: response,
        err: nil,
    }
}

func newErrorOperationResponse(path string, err error) operationResponse {
    return operationResponse {
        path: path,
        response: nil,
        err: err,
    }
}
