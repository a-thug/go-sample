package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	router := SetupRouter()

	resReq := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(resReq, req)

	assert.Equal(t, http.StatusOK, resReq.Code)
	assert.Equal(t, "This is root.", resReq.Body.String())
}

func TestHello(t *testing.T) {
	router := SetupRouter()

	token, _ := newJWT("secret_key", "admin", "admin")
	resReq := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/hello", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	router.ServeHTTP(resReq, req)

	assert.Equal(t, http.StatusOK, resReq.Code)
	var helloRes helloResponse
	if err := json.NewDecoder(resReq.Body).Decode(&helloRes); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "Hello World.", helloRes.Text)
	assert.Equal(t, "admin", helloRes.UserID)

}

func TestRestricted(t *testing.T) {
	router := SetupRouter()

	token, _ := newJWT("secret_key", "someone", "admin")
	resReq := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/hello", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	router.ServeHTTP(resReq, req)

	assert.Equal(t, http.StatusForbidden, resReq.Code)
}

func newJWT(secret string, username string, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(secret))
}

type helloResponse struct {
	Claims struct {
		Exp      int    `json:"exp"`
		Password string `json:"password"`
		UserID   string `json:"user_id"`
	} `json:"claims"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}
