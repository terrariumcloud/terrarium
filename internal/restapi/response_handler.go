package restapi

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/errgo.v2/errors"
)

type ResponseHandler struct{}

func (t *ResponseHandler) Write(rw http.ResponseWriter, data interface{}, statusCode int) {
	resp := &DataResponse{
		Code: statusCode,
		Data: data,
	}
	jsonData, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("+%v", errors.Wrap(err))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	if statusCode != http.StatusNoContent {
		_, _ = rw.Write(jsonData)
	}
}

func (t *ResponseHandler) WriteRaw(rw http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("+%v", errors.Wrap(err))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	if statusCode != http.StatusNoContent {
		_, _ = rw.Write(jsonData)
	}
}

func (t *ResponseHandler) Redirect(rw http.ResponseWriter, r *http.Request, uri string) {
	http.Redirect(rw, r, uri, http.StatusFound)
}
