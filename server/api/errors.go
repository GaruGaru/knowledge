package api

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"net/http"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

func httpErr(w http.ResponseWriter, err error, status int) {
	log.Error(err)
	w.WriteHeader(status)
	response := Error{Message: err.Error()}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
