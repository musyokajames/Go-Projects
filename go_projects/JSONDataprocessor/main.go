package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func JSONOperator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var data []map[string]interface{}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}
}

func main() {

}
