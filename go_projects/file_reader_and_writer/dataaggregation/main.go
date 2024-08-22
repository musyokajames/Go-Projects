package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Create the files
	// open the files
	numFiles := 5
	var aggregatedData []map[string]interface{}

	for i := 0; i < numFiles; i++ {
		fileName := fmt.Sprintf("data%d.json", i)
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error opening file", err)
			continue
		}

		// Read through the files and aggregate the data

		var data []map[string]interface{}
		if err := json.NewDecoder(file).Decode(&data); err != nil {
			fmt.Println("Error Decoding files", fileName, ":", err)
			file.Close()
			continue
		}
		file.Close()

		// Write the aggregated data to a newfile

		aggregatedData = append(aggregatedData, data...)

		outputFile, err := os.Create("aggregated_data.json")
		if err != nil {
			fmt.Println("Error opening file", err)
			return
		}
		defer outputFile.Close()

		if err := json.NewEncoder(outputFile).Encode(aggregatedData); err != nil {
			fmt.Println("Error Encoding files", err)
		}
	}

}
