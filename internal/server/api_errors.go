package server

type errorResponse struct {
	Messages []errorMessage `json:"errors"`
	Status   int            `json:"status"`
}

type errorMessage struct {
	Message string `json:"message"`
}

func newError(code int, messages ...string) *errorResponse {
	res := errorResponse{
		Status:   code,
		Messages: []errorMessage{},
	}

	for _, m := range messages {
		res.Messages = append(res.Messages, errorMessage{Message: m})
	}

	return &res
}
