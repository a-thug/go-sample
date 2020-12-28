package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRoot(t *testing.T) {
	logger := zap.NewExample().Sugar()
	router := http.NewServeMux()
	Setup(router, logger)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	router.ServeHTTP(resRec, req)
	require.Equal(t, http.StatusOK, resRec.Code)            // Check the status code.
	require.Equal(t, "This is root.", resRec.Body.String()) // Check the response body.
}
