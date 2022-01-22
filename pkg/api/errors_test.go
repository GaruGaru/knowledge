package api

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Errors(t *testing.T) {
	r := httptest.NewRecorder()
	const errorMessage = "internal error"
	httpErr(r, errors.New(errorMessage), http.StatusInternalServerError)
	require.Equal(t, http.StatusInternalServerError, r.Code)
	var apiError Error
	err := json.NewDecoder(r.Body).Decode(&apiError)
	require.NoError(t, err)
	require.Equal(t, errorMessage, apiError.Message)
}
