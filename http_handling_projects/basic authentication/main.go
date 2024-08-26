package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Basic") {

		//Remove "Basic" Prefix
		encodedCredentials := strings.TrimPrefix(authHeader, "Basic")

		//Decode the Base64 encoded credentials
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		//Convert decoded bytes to string
		decodedCredentials := string(decodedBytes)

		//Split the decoded credentials into username and password
		parts := strings.SplitN(decodedCredentials, ":", 2)
		if len(parts) != 2 {
			http.Error(w, "401 Unauthorized Access", http.StatusUnauthorized)
			return
		}

		username := parts[0]
		password := parts[1]

		if username == "yourUsername" && password == "yourPassword" {
			fmt.Fprintln(w, "Hello, authenticated user!")
		} else {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		}

	} else {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/secure-endpoint", authenticateUser)

	fmt.Println("Server starting on port 8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
