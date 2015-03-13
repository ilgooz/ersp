package erresp

import (
	"encoding/json"
	"net/http"
)

// add Must() ?

var (
	UnvalidFormData = "Unvalid Form Data"
	BadFormEncoding = "Bad Form Encoding"
)

type Response struct {
	Error Error `json:"error,omitempty"`
	w     http.ResponseWriter
}

type Error struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func New(w http.ResponseWriter) *Response {
	return &Response{
		w:     w,
		Error: Error{Fields: make(map[string]string, 0)},
	}
}

func (er *Response) Message(message string) *Response {
	er.Error.Message = message
	return er
}

func (er *Response) Field(name, message string) *Response {
	er.Error.Fields[name] = message
	return er
}

func (er *Response) Send(status int) error {
	if er.Error.Message == "" && len(er.Error.Fields) > 0 {
		er.Message(UnvalidFormData)
	}

	var data interface{}
	if !er.HasError() {
		data = nil
	} else {
		data = er.Error
	}

	return er.response(status, data)
}

func (er *Response) SendParseFormError() error {
	return er.Message(BadFormEncoding).Send(400)
}

func (er *Response) SendMessage(message string, status int) error {
	return er.Message(message).Send(status)
}

func (er *Response) HasError() bool {
	return er.Error.Message != "" || len(er.Error.Fields) > 0
}

func (er *Response) response(status int, data interface{}) error {
	er.w.WriteHeader(status)
	er.w.Header().Set("Content-type", "application/json; charset=utf-8")
	return json.NewEncoder(er.w).Encode(data)
}
