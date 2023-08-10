package helper

import (
	"encoding/json"
	"net/http"

	"github.com/iruldev/mini-wallet/src/constant"
)

type Response interface {
	ReplyCustom(httpStatusCode int, res any) error
}

type ResponseX struct {
	w      http.ResponseWriter
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func NewResponse(Status string, Data ...any) Response {
	rx := ResponseX{
		Status: Status,
	}
	if len(Data) > 0 {
		rx.Data = Data[0]
	}
	return &rx
}

func PlugResponse(w http.ResponseWriter) Response {
	res := &ResponseX{
		Status: constant.FAILED,
		Data:   nil,
	}
	res.w = w
	return res
}

func (r *ResponseX) ReplyCustom(httpStatusCode int, res any) error {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(httpStatusCode)
	return json.NewEncoder(r.w).Encode(res)
}
