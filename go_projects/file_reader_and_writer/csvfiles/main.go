package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func csvFileOperations(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data.csv")
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		fmt.Println("Error opening file", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error encoding CSV:", err)
		return
	}

	var filteredRecords [][]string

	for _, record := range records {
		name := record[0]
		age := record[1]
		email := record[2]
		//occupation := record[3]

		ageInt, err := strconv.Atoi(age)
		if err != nil {
			fmt.Println("Error converting age:", err)
			continue
		}

		if ageInt > 30 {
			fmt.Println("Name:", name, "Age:", age, "Email:", email)
			filteredRecords = append(filteredRecords, record)
		}
	}

	outputFile, err := os.Create("filtered_data.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)

	header := []string{"Name", "Age", "Email"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header to CSV:", err)
		return
	}

	if err := writer.WriteAll(filteredRecords); err != nil {
		fmt.Println("Error writing records to CSV:", err)
		return
	}

	writer.Flush()

	fmt.Println("Filtered data written to filtered_data_csv")

}

func main() {
	http.HandleFunc("/", csvFileOperations)
	fmt.Println("Starting server on port 9200..")
	http.ListenAndServe(":9200", nil)
}
