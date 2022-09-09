package restapi

import "net/http"

type DataResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ServerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Responder interface {
	Write(rw http.ResponseWriter, data interface{}, statusCode int)
	WriteRaw(rw http.ResponseWriter, data interface{}, statusCode int)
	Redirect(rw http.ResponseWriter, r *http.Request, uri string)
}

type ErrorResponder interface {
	Write(rw http.ResponseWriter, err error, statusCode int)
}
