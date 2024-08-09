package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func binaryFileOperations(w http.ResponseWriter, r *http.Request) {
	//Open the binary file
	inputFile, err := os.Open("binary_input.bin")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	//Create the output binary
	outputFile, err := os.Create("binary_output.bin")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	//Copy the contents from input to output file
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Println("Error copying data:", err)
	}

}

func main() {
	http.HandleFunc("/", binaryFileOperations)
	fmt.Println("Server starting on port 9400...")
	http.ListenAndServe(":9400", nil)
}
