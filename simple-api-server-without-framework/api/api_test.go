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

func TestFoo(t *testing.T) {
	logger := zap.NewExample().Sugar()
	router := http.NewServeMux()
	Setup(router, logger)

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("message", "test")
	req.URL.RawQuery = q.Encode()

	resRec := httptest.NewRecorder()
	router.ServeHTTP(resRec, req)
	require.Equal(t, http.StatusOK, resRec.Code)                    // Check the status code.
	require.Equal(t, "Your message is test.", resRec.Body.String()) // Check the response body.
}

func TestJSON(t *testing.T) {
	logger := zap.NewExample().Sugar()
	router := http.NewServeMux()
	Setup(router, logger)

	req, err := http.NewRequest("GET", "/json", nil)
	if err != nil {
		t.Fatal(err)
	}

	resRec := httptest.NewRecorder()
	router.ServeHTTP(resRec, req)
	require.Equal(t, http.StatusOK, resRec.Code)          // Check the status code.
	require.Equal(t, `{"ok":true}`, resRec.Body.String()) // Check the response body.
}
