package api

import (
	"context"
	"github.com/garugaru/knowledge/server/data"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Router(t *testing.T) {
	api := New(Config{}, nil)
	router := api.router()

	r := httptest.NewRecorder()
	router.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	require.Equal(t, http.StatusOK, r.Code)
}

func TestAPI_RouterForCatalog(t *testing.T) {
	catalog := mockCatalog{
		getDocument: func(ctx context.Context, request data.GetDocumentRequest) (data.Document, error) {
			return data.Document{}, nil
		},
	}

	api := New(Config{}, catalog)
	router := api.router()

	r := httptest.NewRecorder()
	router.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/catalog/documents/1", nil))

	require.Equal(t, http.StatusOK, r.Code)
}
