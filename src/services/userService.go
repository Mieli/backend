package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	pkgmodels "github.com/mieli/backend/src/models"
	pkgrepositories "github.com/mieli/backend/src/repositories"
)

type UserService struct {
	Repository *pkgrepositories.UserRepository
}

func NewUserService(repository *pkgrepositories.UserRepository) *UserService {
	return &UserService{
		Repository: repository,
	}
}

func (s *UserService) FindAll() ([]pkgmodels.UserPresenter, error) {
	users, err := s.Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var listUsers []pkgmodels.UserPresenter
	for _, user := range users {
		listUsers = append(listUsers, pkgmodels.UserPresenter{
			ID:    user.ID,
			Email: user.Email,
		})
	}

	return listUsers, nil
}

func (s *UserService) FindById(id string) (*pkgmodels.UserPresenter, error) {
	user, err := s.Repository.FindById(id)
	if err != nil {
		return nil, err
	}

	return &pkgmodels.UserPresenter{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *UserService) Add(user pkgmodels.User) error {

	user.Email = strings.ToLower(user.Email)

	hashPassword, err := GenerateHashByPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	err = s.Repository.Add(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Update(id string, updatedUser pkgmodels.User) error {
	hashPassword, err := GenerateHashByPassword(updatedUser.Password)
	if err != nil {
		return err
	}

	updatedUser.Password = hashPassword

	err = s.Repository.Update(id, updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Remove(id string) error {
	err := s.Repository.Remove(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindByEmail(email string) (*pkgmodels.User, error) {
	user, err := s.Repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ComparePasswords(user *pkgmodels.User, password string) (bool, error) {
	if user != nil {
		ok := VerifyPasswords(user.Password, password)
		if !ok {
			return false, fmt.Errorf("senha não confere")
		}
		return true, nil
	}
	return false, fmt.Errorf("credenciais inválidas")
}

func GenerateHashByPassword(password string) (string, error) {
	pass := strings.TrimSpace(password)
	hasher := sha256.New()
	_, err := hasher.Write([]byte(pass))
	if err != nil {
		return "", err
	}
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword, nil

}

func VerifyPasswords(passwordDB, password string) bool {

	// Gerar o hash SHA-256 da senha fornecida pelo usuário
	hashedUserPassword, err := GenerateHashByPassword(password)
	if err != nil {
		fmt.Println("Erro ao gerar o hash da senha fornecida:", err)
		return false
	}

	// Comparar o hash gerado com a senha armazenada no banco de dados
	return passwordDB == hashedUserPassword
}
