package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Country struct {
	Name       Name
	Capital    []string
	Region     string
	Landocked  bool
	Population int
}

type Name struct {
	Common   string
	Official string
}

func fetchData(w http.ResponseWriter, r *http.Request) {

	url := "https://restcountries.com/v3.1/all"

	resp, err := http.Get(url)

	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Error: Status code %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var countries []Country
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(countries); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", fetchData)
	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error Starting server:", err)
	}
}
