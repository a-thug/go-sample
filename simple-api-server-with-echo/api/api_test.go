package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	userJSON := `{"name": "hogetest","email": "hoge@example.com","password": "password"}`

	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resRec := httptest.NewRecorder()

	uc := newUserController("hoge")
	c := e.NewContext(req, resRec)

	if assert.NoError(t, uc.Create(c)) {
		assert.Equal(t, http.StatusCreated, resRec.Code)
	}
}
