package utils

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Dev     string      `json:"developer"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(message string, data interface{}) *Response {
	return &Response{
		Status:  "success",
		Message: message,
		Dev:     "Achmad Nashruddin Riskynanda",
		Data:    data,
	}
}

func ResponseError(message string) *Response {
	return &Response{
		Status:  "error",
		Message: message,
		Dev:     "Achmad Nashruddin Riskynanda",
		Data:    nil,
	}

}
