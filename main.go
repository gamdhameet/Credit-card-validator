package main

import (
	"database/sql"
	"log"
	"os"

	//"database/sql"
	"fmt"
	"html/template"
	//"log"
	"net/http"
	//"os"
	"strings"
	//"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type ValidationResult struct {
	Valid      bool
	CardNumber string
	CardType   string
	// Add more fields if needed
}

type Response struct {
	Valid bool `json:"valid"`
}

func main() {
	//tmpl := template.Must(template.ParseFiles("index.html"))
	// Register the creditCardValidator function to handle requests at the root ("/") path.
	//http.HandleFunc("/", creditCardValidator)
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/validate", serveIndex)
	http.HandleFunc("/result", creditCardValidator)
	fmt.Println("Listening on port: 8080")

	// Start an HTTP server listening on the specified port.
	err = http.ListenAndServe(":8080", nil)
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

	err := request.ParseForm()
	if err != nil {
		return
	}
	cardNumber := request.FormValue("cardNumber")
	// Validate the credit card number using the Luhn algorithm.
	isValid := luhnAlgorithm(cardNumber)
	cardType := getCardType(cardNumber)

	result := ValidationResult{
		Valid:      isValid,
		CardNumber: cardNumber,
		CardType:   cardType,
	}
	err = saveCardInfo(result)
	if err != nil {
		http.Error(writer, "Failed to save data", http.StatusInternalServerError)
		log.Printf("Failed to save data: %v", err)
		return
	}
	tmpl2 := template.Must(template.ParseFiles("result.html"))
	tmpl2.Execute(writer, result)
}

func getCardType(cardNumber string) string {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	switch cardNumber[0] {
	case '3':
		return "American Express"
	case '4':
		return "Visa"
	case '5':
		return "MasterCard"
	case '6':
		return "Discover"
	default:
		return "Unknown"
	}
}
func saveCardInfo(result ValidationResult) error {
	query := `
		INSERT INTO card_info (card_number, card_type, cardholder_name, expiry_date)
		VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, result.CardNumber, result.CardType)
	if err != nil {
		log.Printf("Error saving card info: %v\n", err)
		return err
	}
	return nil
}
