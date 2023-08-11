package helper

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request struct {
	r      http.Request
	params map[string]interface{}
	header http.Header
	Method string
}

func PlugRequest(r *http.Request, w http.ResponseWriter) *Request {
	req := &Request{
		r:      *r,
		params: map[string]interface{}{},
		header: r.Header,
		Method: r.Method,
	}

	for k, v := range r.URL.Query() {
		req.params[k] = scan(v)
	}

	switch r.Method {
	case http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch:
		{
			contentType := req.header.Get("Content-Type")
			if strings.Contains(contentType, "multipart/form-data") {
				if r.Method == http.MethodGet {
					return req
				}
				err := r.ParseMultipartForm(32 << 10)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return req
				}
				for k, v := range r.MultipartForm.Value {
					req.params[k] = scan(v)
				}
			}
			break
		}
	}
	return req
}

func scan(values []string) interface{} {
	if len(values) == 1 {
		return identify(values[0])
	} else if len(values) > 1 {
		list := []interface{}{}
		for k, vs := range values {
			list[k] = identify(vs)
		}
		return list
	} else {
		return nil
	}
}

func identify(value string) interface{} {
	var arr []interface{}
	var mp map[string]interface{}
	errArr := json.Unmarshal([]byte(value), &arr)
	errMp := json.Unmarshal([]byte(value), &mp)
	if errArr == nil {
		return arr
	} else if errMp == nil {
		return mp
	} else {
		return value
	}
}

func ParseTo[T any](r *Request) (T, error) {
	jsonString, _ := json.Marshal(r.params)
	var en T
	err := json.Unmarshal(jsonString, &en)
	if err != nil {
		return en, err
	} else {
		return en, nil
	}
}
