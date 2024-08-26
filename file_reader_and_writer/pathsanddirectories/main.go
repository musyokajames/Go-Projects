package main

import (
	"fmt"
	"os"
)

func main() {

	err := os.Mkdir("test_files", 0755)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory created succesfully!")
	}

	// Create some files in the directory
	filesToCreate := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, fileName := range filesToCreate {
		filePath := "test_files/" + fileName
		_, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	}

	outputFile, err := os.Create("file_list.txt")
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	defer outputFile.Close()

	files, err := os.ReadDir("./test_files")
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, file := range files {
		if _, err := outputFile.WriteString(file.Name() + "\n"); err != nil {
			fmt.Println("Error writing to file_list.txt:", err)
			return
		}
	}
	fmt.Println("File names written to file_list.txt succesfully")
}
