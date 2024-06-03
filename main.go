package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var client *mongo.Client

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")
}

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
		//http.Error(writer, "Failed to save data", http.StatusInternalServerError)
		log.Printf("Failed to save data: %v", err)
	}
	tmpl2 := template.Must(template.ParseFiles("result.html"))
	tmpl2.Execute(writer, result)
	//log.Printf("Failed to save data: %v", err)
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

type CardInfo struct {
	Valid      string `bson:"valid"`
	CardNumber string `bson:"cardNumber"`
	CardType   string `bson:"cardType"`
}

func saveCardInfo(result ValidationResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database("your_database").Collection("card_info")
	card := CardInfo{
		Valid:      fmt.Sprintf("%t", result.Valid),
		CardNumber: result.CardNumber,
		CardType:   result.CardType,
	}
	_, err := collection.InsertOne(ctx, card)
	return err
}
