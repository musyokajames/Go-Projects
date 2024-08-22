package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var tasks []Task

func addTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}
	tasks = append(tasks, newTask)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func viewTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	idStr := params["id"]

	// Convert the id from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for index, item := range tasks {
		if item.ID == id {
			tasks = append(tasks[:index], tasks[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	tasks = append(tasks, Task{ID: 1, Task: "Wake up"})
	tasks = append(tasks, Task{ID: 2, Task: "Gym"})
	tasks = append(tasks, Task{ID: 3, Task: "Shower"})
	tasks = append(tasks, Task{ID: 4, Task: "Work"})

	r.HandleFunc("/tasks", addTasks).Methods("POST")
	r.HandleFunc("/tasks", viewTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", deleteTasks).Methods("DELETE")

	fmt.Printf("Starting server on port 8000...\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
