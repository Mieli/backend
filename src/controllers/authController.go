package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	pkgmodels "github.com/mieli/backend/src/models"
	pkgservice "github.com/mieli/backend/src/services"
)

type AuthController struct {
	Service   *pkgservice.UserService
	SecretKey string
	Token     string
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthController(service *pkgservice.UserService) *AuthController {
	return &AuthController{
		Service:   service,
		SecretKey: os.Getenv("SECRET_KEY"),
		Token:     "",
	}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.Service.FindByEmail(credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := ac.Service.ComparePasswords(user, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok {

		// criar o token
		ac.Token, _ = GenerateTokenJWT(*user, []byte(ac.SecretKey))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ac.Token))

	} else {
		fmt.Println("erro ao verificar a senha")
	}

}

func (ac *AuthController) Logoff(w http.ResponseWriter, r *http.Request) {
	ac.SecretKey = ""
	ac.Token = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Logoff bem-sucedido!")

}

func GenerateTokenJWT(user pkgmodels.User, key []byte) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Erro ao Gerar o Token do usuário")
	}
	return token, nil
}
