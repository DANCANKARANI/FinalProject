package main

import (
	"fmt"
	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/api/routes"
)

func main() {
	model.DbMigrator()
	fmt.Println("hello medicare")
	routes.RegisterEndpoints()
}

/*import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Replace these with your actual credentials




func main() {
	err := godotenv.Load(".env")
    if err != nil {
        fmt.Printf("Error loading .env file: %v", err)
    }
	// Step 1: Generate Access Token
	accessToken, err := generateAccessToken()
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return
	}

	fmt.Println("Access Token:", accessToken)

	// Step 2: Make STK Push Request
	stkPushResponse, err := makeSTKPushRequest(accessToken)
	if err != nil {
		fmt.Println("Error making STK Push request:", err)
		return
	}

	fmt.Println("STK Push Response:", stkPushResponse)
}

// Function to generate access token
func generateAccessToken() (string, error) {
	auth := os.Getenv("Safaricom_ConsumerKey") + ":" + os.Getenv("Safaricom_ConsumerSecret")
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", basicAuth)

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response to extract the access token
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return accessToken, nil
}

// Function to make STK Push request
func makeSTKPushRequest(accessToken string) (string, error) {
	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	method := "POST"

	payloadData := map[string]interface{}{
		"BusinessShortCode": 174379,
		"Password":         os.Getenv("Safaricom_Password"),
		"Timestamp":        "20231217153132",
		"TransactionType":  "CustomerPayBillOnline",
		"Amount":           80,
		"PartyA":           25497408042,
		"PartyB":           174379,
		"PhoneNumber":      254797408042,
		"CallBackURL":      "http://192.168.0.101",
		"AccountReference": "Penta Drive",
		"TransactionDesc":  "Bike Ride",
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}*/