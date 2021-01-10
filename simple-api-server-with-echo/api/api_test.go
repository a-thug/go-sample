package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

var (
	secret      = "test"
	token       = ""
	createdUser = &user{
		ID:           1,
		Name:         "hogetest",
		Email:        "hoge@example.com",
		PasswordHash: []byte(""),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Time{},
	}
)

func TestCreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	userJSON := `{"name": "hogetest","email": "hoge@example.com","password": "password"}`

	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resRec := httptest.NewRecorder()

	uc := newUserController(secret)
	c := e.NewContext(req, resRec)

	if assert.NoError(t, uc.Create(c)) {
		assert.Equal(t, http.StatusCreated, resRec.Code)
		var createRes userCreateResponse
		if err := json.Unmarshal(resRec.Body.Bytes(), &createRes); err != nil {
			log.Fatal(err)
		}
		token = createRes.Token
		assert.Equal(t, int64(1), createRes.User.ID)
		assert.Equal(t, "hogetest", createRes.User.Name)
		assert.Equal(t, "hoge@example.com", createRes.User.Email)
	}

}

func TestGetUser(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	resRec := httptest.NewRecorder()
	uc := newUserController(secret)
	uc.users[1] = createdUser

	c := e.NewContext(req, resRec)
	c.SetPath("/api/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	exec := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
	})(uc.Get)(c)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, resRec.Code)
		var getRes userGetResponse
		if err := json.Unmarshal(resRec.Body.Bytes(), &getRes); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, int64(1), getRes.User.ID)
		assert.Equal(t, "hogetest", getRes.User.Name)
		assert.Equal(t, "hoge@example.com", getRes.User.Email)
	}
}

func TestUpdateUser(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}

	userJSON := `{"name": "testtest","email": "hogehoge@example.com"}`
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resRec := httptest.NewRecorder()
	uc := newUserController(secret)
	uc.users[1] = createdUser

	c := e.NewContext(req, resRec)
	c.SetPath("/api/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	exec := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
	})(uc.Update)(c)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusOK, resRec.Code)
		var updateRes userUpdateResponse
		if err := json.Unmarshal(resRec.Body.Bytes(), &updateRes); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, int64(1), updateRes.User.ID)
		assert.Equal(t, "testtest", updateRes.User.Name)
		assert.Equal(t, "hogehoge@example.com", updateRes.User.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	resRec := httptest.NewRecorder()
	uc := newUserController(secret)
	uc.users[1] = createdUser

	c := e.NewContext(req, resRec)
	c.SetPath("/api/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	exec := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
	})(uc.Delete)(c)

	if assert.NoError(t, exec) {
		assert.Equal(t, http.StatusNoContent, resRec.Code)
		assert.Nil(t, uc.users[1])
	}
}
