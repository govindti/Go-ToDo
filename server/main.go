package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var (
	todos     []Todo
	todoMutex sync.Mutex
)

type HeathResponse struct {
	Status  string `json:"status`
	Message string `json:"message`
}

type Todo struct {
	Id        string `json:"id`
	Task      string `json:"task`
	Completed bool   `json:"completed`
}

type TodobyID struct {
	Status  string `json:"status`
	Message string `json:"message`
}

func genrateRandomId() string {
	return uuid.New().String()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HeathResponse{
		Status:  "Ok",
		Message: "Api is Up & Running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)

	case "POST":
		var newTodo Todo
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Unable to read from req body", http.StatusBadRequest)
		}
		err = json.Unmarshal(body, &newTodo)

		if err != nil || newTodo.Task == "" {
			http.Error(w, "No Input", http.StatusBadRequest)
			return
		}
		newTodo.Id = genrateRandomId()
		todoMutex.Lock()
		todos = append(todos, newTodo)
		todoMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}

}

func todoByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	todoMutex.Lock()
	defer todoMutex.Unlock()

	for i, todo := range todos {
		if todo.Id == id {
			switch r.Method {
			case "GET":
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(todo)

			case "PUT":
				var updatedTodo Todo
				body, err := ioutil.ReadAll(r.Body)

				if err != nil {
					http.Error(w, "Unable to read from req body", http.StatusBadRequest)
					return
				}

				err = json.Unmarshal(body, &updatedTodo)
				if err != nil || updatedTodo.Task == "" {
					http.Error(w, "Invalid Input", http.StatusBadRequest)
					return
				}

				todos[i].Task = updatedTodo.Task
				todos[i].Completed = updatedTodo.Completed

				json.NewEncoder(w).Encode(todos[i])

			case "DELETE":
				deletedTodoId := todo.Id
				todos = append(todos[:i], todos[i+1:]...)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Todo Deleted with id " + deletedTodoId + " Succesfully"})
			default:
				http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)

			}
		}
	}
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoByIdHandler)

	fmt.Println("App is running on PORT 3000")

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error Starting the App", err)
	}

}
