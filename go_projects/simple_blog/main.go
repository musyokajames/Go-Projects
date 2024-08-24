package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var posts []Post

func viewPosts(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, posts)
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	for _, item := range posts {
		if item.ID == id {
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, item)
			return
		}
	}
	http.Error(w, "Post not found", http.StatusNotFound)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		newPost := Post{
			ID:      len(posts) + 1,
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
		}

		posts = append(posts, newPost)

		http.Redirect(w, r, "/post/", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "templates/new_post.html")
		return
	}

}

func deletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for index, item := range posts {
		if item.ID == id {
			posts = append(posts[:index], posts[index+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/post/", http.StatusSeeOther)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/post/", viewPosts).Methods("GET")
	r.HandleFunc("/post/{id}", viewPost).Methods("GET")
	r.HandleFunc("/new/", createPost).Methods("POST", "GET")
	r.HandleFunc("/delete/{id}", deletePost).Methods("POST")

	fmt.Println("Server starting at port 9000....")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatal(err)
	}
}
