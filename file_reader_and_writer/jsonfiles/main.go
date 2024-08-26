package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Person struct {
	Id           int
	Name         string
	Age          int
	Email        string
	isActive     bool
	Address      Address
	PhoneNumbers []string
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func jsonFileOperations(w http.ResponseWriter, r *http.Request) {
	//open the JSON file for reading
	file, err := os.Open("data.json")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var persons []Person

	//Decode JSON data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&persons)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	//Modify data
	for i := range persons {
		persons[i].Email = "updated_" + persons[i].Email
	}

	//Create a new file to write the updated data
	outputFile, err := os.Create("updated_data.json")
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	//Encode updataed data to JSON and write to the new file
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", " ") //format the JSON output
	err = encoder.Encode(persons)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Fprintln(w, "JSON file processed and saved as updated_data.json")

}

func main() {
	http.HandleFunc("/", jsonFileOperations)
	fmt.Println("Server starting on port 9300...")
	http.ListenAndServe(":9300", nil)
}
