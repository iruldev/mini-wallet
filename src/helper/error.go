package helper

import (
	"github.com/go-playground/validator/v10"
)

type ErrField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrMsgField(err error) (errF []ErrField) {
	for _, e := range err.(validator.ValidationErrors) {
		errF = append(errF, ErrField{
			Field:   e.Field(),
			Message: e.Tag(),
		})
	}
	return
}
