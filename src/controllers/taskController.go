package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	pkgmodels "github.com/mieli/backend/src/models"
	pkgservices "github.com/mieli/backend/src/services"
)

type TaskController struct {
	Service *pkgservices.TaskService
}

func NewTaskController(service *pkgservices.TaskService) *TaskController {
	return &TaskController{
		Service: service,
	}
}

func (c *TaskController) FindAll(w http.ResponseWriter, r *http.Request) {
	tasks, _ := c.Service.FindAll()
	response, _ := json.Marshal(tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (c *TaskController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["id"]
	task, _ := c.Service.FindById(taskId)
	response, _ := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (c *TaskController) Add(w http.ResponseWriter, r *http.Request) {

	var task pkgmodels.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.Service.Add(task)
	if err != nil {
		http.Error(w, "Erro ao inserir a tarefa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (c *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["id"]

	var updatedTask pkgmodels.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.Service.Update(taskId, updatedTask)
	if err != nil {
		http.Error(w, "Erro ao remover a tarefa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *TaskController) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["id"]

	err := c.Service.Remove(taskId)
	if err != nil {
		http.Error(w, "Erro ao remover a tarefa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
