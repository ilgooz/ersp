package ersp

import (
	"encoding/json"
	"net/http"
	"strings"
)

// add Must() ?

var (
	UnvalidFormData  = "Unvalid Form Data"
	BadFormEncoding  = "Bad Form Encoding"
	PrettyErrorTitle = "Error"
)

var (
	MustBeANumberErr               = "must be a number"
	MustBeNumbersErr               = "must be numbers"
	MustBeCommaSeparatedNumbersErr = "must be comma separated numbers"
	MustBeTrueOrFalseErr           = "must be true or false"
	MustBeUTCErr                   = "must be UTC"
)

type Response struct {
	Error Error `json:"error,omitempty"`
	w     http.ResponseWriter
	r     *http.Request
}

type Error struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func New(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{
		w:     w,
		r:     r,
		Error: Error{Fields: make(map[string]string, 0)},
	}
}

func (er *Response) isPretty() bool {
	return er.r.Header.Get("X-Pretty-Error") == "true"
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
		if er.isPretty() {
			er.Message(PrettyErrorTitle)
		} else {
			er.Message(UnvalidFormData)
		}
	}

	if er.isPretty() {
		fields := map[string]string{}

		for key, val := range er.Error.Fields {
			newKey := strings.Title(strings.Replace(key, "_", " ", -1))

			var newVal string
			switch val {
			case MustBeANumberErr, MustBeNumbersErr, MustBeCommaSeparatedNumbersErr,
				MustBeTrueOrFalseErr, MustBeUTCErr:
				newVal = "must be specified"
			default:
				newVal = val
			}

			fields[newKey] = newVal
		}

		er.Error.Fields = fields
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
	er.w.Header().Set("Content-type", "application/json; charset=utf-8")
	er.w.WriteHeader(status)
	return json.NewEncoder(er.w).Encode(data)
}
