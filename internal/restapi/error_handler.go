package restapi

import (
	"encoding/json"
	"fmt"
	"gopkg.in/errgo.v2/errors"
	"log"
	"net/http"
)

const BadRequestPrefix string = "Bad Request"
const InternalServerErrorPrefix string = "Internal Server Error"
const NotFoundPrefix string = "404 Not Found"
const UnprocessablePrefix string = "Unprocessable Entity"
const NotImplementedPrefix string = "Not Implemented"

type ErrorHandler struct{}

func (t *ErrorHandler) Write(rw http.ResponseWriter, err error, statusCode int) {
	var prefix string = ""
	switch statusCode {
	case http.StatusInternalServerError:
		prefix = InternalServerErrorPrefix
	case http.StatusBadRequest:
		prefix = BadRequestPrefix
	case http.StatusNotFound:
		prefix = NotFoundPrefix
	case http.StatusUnprocessableEntity:
		prefix = UnprocessablePrefix
	case http.StatusNotImplemented:
		prefix = NotImplementedPrefix
	default:

	}
	resp := &ServerResponse{
		Code:    statusCode,
		Message: fmt.Sprintf("%s - %s", prefix, err.Error()),
	}
	jsonData, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("+%v", errors.Wrap(err))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	_, _ = rw.Write(jsonData)
}
