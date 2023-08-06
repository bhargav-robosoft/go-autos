package service

import (
	"autos/db"
	"autos/middleware"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	Login(email string, password string) (token string, err error)
	Register(email string, password string) (err error)
	Logout(token string)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (service *authService) Login(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("Email and password are required")
	}

	validCredentials, userId, err := db.CheckUserCredentials(email, password)
	if err != nil {
		return "", err
	}

	if !validCredentials {
		return "", errors.New("Email or password is wrong")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["id"] = userId

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (service *authService) Register(email string, password string) error {
	if email == "" || password == "" {
		return errors.New("Email and password are required")
	}

	isNewUser, err := db.CheckNewUser(email)
	if err != nil {
		return err
	}

	if !isNewUser {
		return errors.New("User with this email is already registered")
	}

	_, err = db.RegisterUser(email, password)
	if err != nil {
		return err
	}

	return nil
}

func (service *authService) Logout(token string) {
	middleware.BlackListTokens = append(middleware.BlackListTokens, token)
}
