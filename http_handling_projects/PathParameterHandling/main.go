package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func greetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Fprintf(w, "Hello, %s", name)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/greet/{name:[a-zA-Z]+}", greetHandler)

	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
