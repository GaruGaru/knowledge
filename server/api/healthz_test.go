package api

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Healthz(t *testing.T) {
	api := New(Config{}, nil)
	r := httptest.NewRecorder()
	api.healthz(r, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	require.Equal(t, http.StatusOK, r.Code)
}
