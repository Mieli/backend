package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	pkgcontrollers "github.com/mieli/backend/src/controllers"
	pkgdatabase "github.com/mieli/backend/src/database"
	pkgrepositories "github.com/mieli/backend/src/repositories"
	pkgservices "github.com/mieli/backend/src/services"
)

type Dependecies struct {
	TaskService *pkgservices.TaskService
	UserService *pkgservices.UserService
}

func createDependecies() Dependecies {

	// criar a instância do banco de dados
	dbMongo := pkgdatabase.DataBaseMongo{}
	client, err := dbMongo.Connect()
	if err != nil {
		fmt.Println("Erro na conexao com o banco de dados" + err.Error())
	}

	taskRepository := pkgrepositories.NewTaskRepository(client)
	taskService := pkgservices.NewTaskService(taskRepository)

	userRepository := pkgrepositories.NewUserRepository(client)
	userService := pkgservices.NewUserService(userRepository)

	return Dependecies{
		TaskService: taskService,
		UserService: userService,
	}
}

func configRouter(dependencies *Dependecies) *mux.Router {

	mainRouter := mux.NewRouter()

	if dependencies.UserService != nil {
		authController := pkgcontrollers.NewAuthController(dependencies.UserService)
		if authController != nil {
			mainRouter.HandleFunc("/login", authController.Login).Methods("POST")
			mainRouter.HandleFunc("/logoff", authController.Logoff).Methods("GET")
		}
	}

	apiRouter := mainRouter.PathPrefix("/api/v1").Subrouter()

	if dependencies.TaskService != nil && dependencies.UserService != nil {

		taskController := pkgcontrollers.NewTaskController(dependencies.TaskService)
		apiRouter.HandleFunc("/tasks", taskController.FindAll).Methods("GET")
		apiRouter.HandleFunc("/tasks/{id}", taskController.FindById).Methods("GET")
		apiRouter.HandleFunc("/tasks", taskController.Add).Methods("POST")
		apiRouter.HandleFunc("/tasks/{id}", taskController.Update).Methods("PUT")
		apiRouter.HandleFunc("/tasks/{id}", taskController.Remove).Methods("DELETE")

		userController := pkgcontrollers.NewUserController(dependencies.UserService)
		apiRouter.HandleFunc("/users", userController.FindAll).Methods("GET")
		apiRouter.HandleFunc("/users/{id}", userController.FindById).Methods("GET")
		apiRouter.HandleFunc("/users", userController.Add).Methods("POST")
		apiRouter.HandleFunc("/users/{id}", userController.Update).Methods("PUT")
		apiRouter.HandleFunc("/users/{id}", userController.Remove).Methods("DELETE")
	}

	return mainRouter
}

func readDotEnv() {
	// Carrega as variáveis de ambiente de um arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
}

func Start() {
	readDotEnv()

	dependencies := createDependecies()
	router := configRouter(&dependencies)

	port := 3000

	fmt.Println("Aplicação rodando http://localhost:" + strconv.Itoa(port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Println(err)
	}

}
