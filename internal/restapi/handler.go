package restapi

import "net/http"

type RESTAPIHandler interface {
	GetHttpHandler(mountPath string) http.Handler
}
