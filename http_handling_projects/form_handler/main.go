package main

import (
	"fmt"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Post Request Successful\n")
	Name := r.FormValue("name")
	message := r.FormValue("message")
	fmt.Fprintf(w, "Name: %s\n", Name)
	fmt.Fprintf(w, "Message: %s\n", message)

}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/form", formHandler)
	fmt.Println("Starting server on port 8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
