package api

import (
	"net/http"
)

func (a Api) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
