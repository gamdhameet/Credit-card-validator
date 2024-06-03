# Credit Card Validator

This project is a simple web application that validates credit card numbers using the Luhn algorithm and stores the card details in a MongoDB database.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.13+)
- [MongoDB](https://www.mongodb.com/try/download/community) (running locally or on a remote server)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/credit-card-validator.git
    cd credit-card-validator
    ```

2. Install Go dependencies:

    ```sh
    go mod tidy
    ```

## Setup MongoDB

1. Start your MongoDB server. For a local MongoDB instance, you can use the following command:

    ```sh
    mongod --dbpath /path/to/your/mongodb/data
    ```

2. Ensure your MongoDB server is running and accessible at `mongodb://localhost:27017`.

## Project Structure

- `main.go`: The main application file.
- `index.html`: The HTML form for inputting credit card numbers.
- `result.html`: The HTML page for displaying validation results.

## Run the Project

1. Open a terminal and navigate to the project directory.

2. Run the application:

    ```sh
    go run main.go
    ```

3. Open your web browser and navigate to `http://localhost:8080`.

## Usage

1. Enter a credit card number in the form on the index page.

2. Submit the form to validate the credit card number.

3. The result page will display whether the card is valid and show the card type and number.