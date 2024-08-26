package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func FileOperations(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("file1.txt")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()

	outputFile, err := os.Create("file2.txt")
	if err != nil {
		http.Error(w, "Error Creating file", http.StatusInternalServerError)
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ToUpper(line)

		_, err := outputFile.WriteString(line + "\n")
		if err != nil {
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Fprintln(w, "File processed and saved as file2.txt")

}

func main() {
	http.HandleFunc("/", FileOperations)
	fmt.Println("Server starting on port 8090...")
	http.ListenAndServe(":9000", nil)
}
