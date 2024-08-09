package main

import (
	"fmt"
	"net/http"
)

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		name := ""
		if name == "" {
			fmt.Fprintln(w, "Hello, Stranger!")
		} else {
			fmt.Fprintln(w, "Hello", name)
		}
	} else {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
}
func main() {
	http.HandleFunc("/hello", greetingHandler)
	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
