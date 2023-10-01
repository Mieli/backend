package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	pkgmodels "github.com/mieli/backend/src/models"
	pkgservices "github.com/mieli/backend/src/services"
)

type UserController struct {
	Service *pkgservices.UserService
}

func NewUserController(service *pkgservices.UserService) *UserController {
	return &UserController{
		Service: service,
	}
}

func (c *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	users, _ := c.Service.FindAll()
	response, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserId := vars["id"]
	User, _ := c.Service.FindById(UserId)
	response, _ := json.Marshal(User)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (c *UserController) Add(w http.ResponseWriter, r *http.Request) {

	var User pkgmodels.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.Service.Add(User)
	if err != nil {
		http.Error(w, "Erro ao inserir o usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserId := vars["id"]

	var updatedUser pkgmodels.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.Service.Update(UserId, updatedUser)
	if err != nil {
		http.Error(w, "Erro ao remover o usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UserController) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserId := vars["id"]

	err := c.Service.Remove(UserId)
	if err != nil {
		http.Error(w, "Erro ao remover o usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
