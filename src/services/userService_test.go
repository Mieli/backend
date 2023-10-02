package services_test

import (
	"testing"

	pkgservice "github.com/mieli/backend/src/services"
)

func TestGenerateHashByPassword(t *testing.T) {
	tests := []struct {
		password        string
		expectedHashLen int
	}{
		{"123", 64}, // O hash SHA-256 deve ter 64 caracteres em formato hexadecimal
		{"senha123", 64},
		{"123456", 64},
		{"", 64}, // Senha vazia ainda deve gerar um hash de 64 caracteres
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			hashedPassword, err := pkgservice.GenerateHashByPassword(test.password)
			if err != nil {
				t.Fatalf("Erro ao gerar o hash: %v", err)
			}

			if len(hashedPassword) != test.expectedHashLen {
				t.Errorf("Tamanho do hash incorreto. Esperado: %d, Obtido: %d", test.expectedHashLen, len(hashedPassword))
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	passwordDB := "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4"
	userProvidedPassword := "1234"

	// Teste de senha correta
	match := pkgservice.VerifyPasswords(passwordDB, userProvidedPassword)
	if !match {
		t.Errorf("Esperado correspondência, mas não houve correspondência")
	}

	// Teste de senha incorreta
	incorrectPassword := "senha456" // Senha diferente da armazenada no banco de dados
	noMatch := pkgservice.VerifyPasswords(passwordDB, incorrectPassword)
	if noMatch {
		t.Errorf("Esperado não correspondência, mas houve correspondência")
	}
}
