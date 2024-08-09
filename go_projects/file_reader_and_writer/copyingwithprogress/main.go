package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func CopyWithProgress(w http.ResponseWriter, r *http.Request) {

	srcFile, err := os.Open("source_file.txt")
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	destFile, err := os.Create("destination_file.txt")
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destFile.Close()

	bar := progressbar.DefaultBytes(fileSize, "Copying")

	buffer := make([]byte, 32*1024)
	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading source file:", err)
			return
		}
		if n == 0 {
			break
		}

		_, err = destFile.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to destination file:", err)
			return
		}

		bar.Add(int(n))
	}

	fmt.Println("File copied succesfully")

}
func main() {
	http.HandleFunc("/", CopyWithProgress)
	fmt.Println("Server opening on port 9600...")
	http.ListenAndServe(":9600", nil)
}
