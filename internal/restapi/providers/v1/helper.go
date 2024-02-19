package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProviderNameFromRequest(r *http.Request) string {
	params := mux.Vars(r)
	orgName := params["organization_name"]
	providerName := params["name"]
	return fmt.Sprintf("%s/%s", orgName, providerName)
}

func GetProviderInputsFromRequest(r *http.Request) (string, string, string) {
	params := mux.Vars(r)
	return params["version"], params["os"], params["arch"]
}