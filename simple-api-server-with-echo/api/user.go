package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type userController struct {
	secret         string
	users          map[int64]*user
	userIDCountter int64
}

func newUserController(secret string) *userController {
	return &userController{secret, make(map[int64]*user), 0}
}

type userCreateRequest struct {
	Name     string `json:"name" validate:"required,gte=1"`
	Email    string `json:"email" validate: "required,email"`
	Password string `json:"password" validate:"required, gte=8"`
}

type userCreateResponse struct {
	User  userInfo `json:"user"`
	Token string   `json:"token"`
}

func (uc *userController) Create(c echo.Context) error {
	r := new(userCreateRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// NOT thread safe
	uc.users[u.ID] = u

	token, err := newJWT(uc.secret, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userCreateResponse{
		User: userInfo{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Token: token,
	})
}

type userGetResponse struct {
	User userInfo `json:"user"`
}

func (uc *userController) Get(c echo.Context) error {
	idParam := c.Param("id")

	id, _ := strconv.ParseInt(idParam, 10, 64)

	session := getSession(c)
	if session.ID != id {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "access denied",
		})
	}

	u, ok := uc.users[id]
	if !ok {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": fmt.Sprintf("could not find user id %d", id),
		})
	}

	return c.JSON(http.StatusOK, userGetResponse{
		User: userInfo{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			CreatedAt, u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	})
}

type userUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,gte=1,required_without_all"`
	Email    string `json:"email" validate:"omitempty,email,required_without_all"`
	Password string `json:"password" validate:"omitempty,gte=8,required_without_all"`
}

type userUpdateResponse struct {
	User userInfo `json:"user"`
}

func (uc *userController) Update(c echo.Context) error {
	idParam := c.Param("id")

	id, _ := strconv.ParseInt(idParam, 10, 64)

	session := getSession(c)
	if session.ID != id {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "access denied",
		})
	}

	u, ok := uc.users[id]
	if !ok {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": fmt.Sprintf("could not find user id %d", id),
		})
	}

	r := new(userUpdateRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if r.Name != "" {
		u.Name = r.Name
	}
	if r.Email != "" {
		u.Email = r.Email
	}
	if r.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = hash
	}
	u.UpdatedAt = time.Now()

	// Update
	uc.users[u.ID] = u

	return c.JSON(http.StatusOK, userGetResponse{
		User: userInfo{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	})
}

func newJWT(secret string, u *user) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func getSession(c echo.Context) session {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	id := claims["id"].(float64)
	name := claims["name"].(string)
	email := claims["email"].(string)

	return session{int64(id), name, email}
}

type session struct {
	ID    int64
	Name  string
	Email string
}
