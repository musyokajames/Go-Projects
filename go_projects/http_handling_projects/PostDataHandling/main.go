package main

import (
	"fmt"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		message := r.FormValue("message")

		if message == "" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		} else {
			fmt.Fprintln(w, message)
		}
	} else {
		http.Error(w, "Method Not ALlowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", formHandler)
	fmt.Println("Starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
