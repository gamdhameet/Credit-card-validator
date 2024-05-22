To test if it's working, use Postman or Curl to send a POST request to the server at route "/", adding a credit card number to the request body:

Set the headers: Content-Type: application/json

Send a POST request to http://localhost:8080/

Set the request body to { "number": "4003600000000014" }

If all goes well, you should receive this response {"valid": true}
