package main

import (
	"fmt"
	"net/http"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("/home/musyoka/Desktop/ai_blog_app/FRONTEND"))

	fs.ServeHTTP(w, r)

}

func main() {

	http.HandleFunc("/", FileHandler)
	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
