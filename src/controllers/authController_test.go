package controllers_test

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt"
	pkgcontrollers "github.com/mieli/backend/src/controllers"
	pkgmodels "github.com/mieli/backend/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
	fmt.Println("ok")
}

func TestGenerateToken(t *testing.T) {
	key := []byte("chavesecreta")
	user := pkgmodels.User{
		ID:       primitive.NewObjectID(),
		Email:    "teste@teste.com.br",
		Password: "123",
	}
	tokenString, err := pkgcontrollers.GenerateTokenJWT(user, key)

	// Verifique se não houve erro na geração do token
	if err != nil {
		t.Errorf("Erro inesperado ao gerar token: %v", err)
	}

	// Decodifique o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifique o algoritmo de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		t.Fatalf("Erro ao decodificar token: %v", err)
	}

	// Verifique se as reivindicações do token estão corretas
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Erro ao obter as reivindicações do token")
	}

	// Verifique as reivindicações específicas que deseja testar
	expectedSub := user.ID.Hex()
	actualSub := claims["sub"].(string)
	if actualSub != expectedSub {
		t.Errorf("Reivindicação 'sub' incorreta. Esperada: %s, Obtida: %s", expectedSub, actualSub)
	}

	expectedEmail := user.Email
	actualEmail := claims["email"].(string)
	if actualEmail != expectedEmail {
		t.Errorf("Reivindicação 'email' incorreta. Esperada: %s, Obtida: %s", expectedEmail, actualEmail)
	}
}
