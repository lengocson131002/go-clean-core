package broker

type Result struct {
	Status  int         `json:"-"`       // Http status code
	Code    string      `json:"code"`    // Error code
	Message string      `json:"message"` // Message
	Details interface{} `json:"details"` //
}

type Response[T any] struct {
	Result Result `json:"result"` // Result
	Data   T      `json:"data"`   // Data
}

var (
	DefaultSuccessResponse = Response[interface{}]{
		Result: Result{
			Code:    "0",
			Message: "Success",
		},
		Data: nil,
	}

	DefaultErrorResponse = Response[interface{}]{
		Result: Result{
			Code:    "1",
			Message: "Internal Server Error",
		},
	}
)

func SuccessResponse[T any](data T) Response[T] {
	var defRes = DefaultSuccessResponse
	return Response[T]{
		Result: defRes.Result,
		Data:   data,
	}
}
