package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string `json:"title"`
	Heading string `json:"heading"`
	Message string `json:"message"`
}

func tmplHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template.html"))

	data := PageData{
		Title:   "Welcome Page",
		Heading: "Hello, World!",
		Message: "This is a message from me to you :)",
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", tmplHandler)
	fmt.Println("Server starting on port 8090...")
	http.ListenAndServe(":8090", nil)

}
