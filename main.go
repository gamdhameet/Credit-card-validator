package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Response struct {
	Valid bool `json:"valid"`
}

func main() {
	//tmpl := template.Must(template.ParseFiles("index.html"))
	// Register the creditCardValidator function to handle requests at the root ("/") path.
	http.HandleFunc("/", creditCardValidator)
	http.HandleFunc("/validate", serveIndex)
	fmt.Println("Listening on port: 8080")

	// Start an HTTP server listening on the specified port.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err) // Print an error message if the server fails to start.
	}
}
func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

// creditCardValidator handles the credit card validation logic and JSON response.
func creditCardValidator(writer http.ResponseWriter, request *http.Request) {
	// Check if the request method is POST.

	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	request.ParseForm()
	cardNumber := request.FormValue("cardNumber")
	// Validate the credit card number using the Luhn algorithm.
	isValid := luhnAlgorithm(cardNumber)
	// Create a response struct with the validation result.
	response := Response{Valid: isValid}

	// Marshal the response struct into JSON format.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error creating response", http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate JSON response.
	writer.Header().Set("Content-Type", "application/json")

	// Write the JSON response back to the client.
	writer.Write(jsonResponse)
}
