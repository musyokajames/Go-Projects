package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Person struct {
	Name      string
	Age       int
	isStudent bool
	Courses   []string
	Address   Address
}

type Address struct {
	Street string
	City   string
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		person := Person{
			Name:      "James Musyoka",
			Age:       21,
			isStudent: true,
			Courses:   []string{"Mathematics", "Computer Science"},
			Address: Address{
				Street: "Mtaa Wa Makka",
				City:   "Mombasa",
			},
		}
		//Set content-type header
		w.Header().Set("Content-Type", "application/json")

		//Encode the person struct to JSON and write it to the response
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", JSONHandler)
	fmt.Println("Server startin on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
