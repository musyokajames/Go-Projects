package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func lineByLineOperations(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("lines.txt")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	outputFile, err := os.Create("processed_lines.txt")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file:", err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)
	lineNumber := 1

	for scanner.Scan() {
		line := fmt.Sprintf("%d: %s\n", lineNumber, scanner.Text())
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}

	writer.Flush()
}

func main() {
	http.HandleFunc("/", lineByLineOperations)
	fmt.Println("Server starting on port 9500...")
	http.ListenAndServe(":9500", nil)
}
