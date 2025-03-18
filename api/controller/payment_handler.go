package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dancankarani/medicare/api/model"
	"github.com/dancankarani/medicare/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MakePayments(c *fiber.Ctx)error{
	 /*_ := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }*/
	// Step 1: Generate Access Token
	accessToken, err := generateAccessToken()
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return utilities.ShowError(c,"failed to send an STK push",1,map[string][]string{"errors":{err.Error()}})
	}

	fmt.Println("Access Token:", accessToken)

	// Step 2: Make STK Push Request
	stkPushResponse, err := makeSTKPushRequest(c,accessToken)
	if err != nil {
		fmt.Println("Error making STK Push request:", err)
		return utilities.ShowError(c,"failed to send an STK push",1,map[string][]string{"errors":{err.Error()}})
	}

	fmt.Println("STK Push Response:", stkPushResponse)
	return utilities.ShowSuccess(c,"Check your mobile phone for an MPESA STK push",0,stkPushResponse)

}


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
func makeSTKPushRequest(c *fiber.Ctx, accessToken string) (*model.Payments, error) {
	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	method := "POST"

	payment := model.Payments{}

	// Parse the payment data from the request body
	if err := c.BodyParser(&payment); err != nil {
		return nil, fmt.Errorf("failed to parse payment JSON data: %v", err)
	}

	// Prepare the payload for the STK Push request
	payloadData := map[string]interface{}{
		"BusinessShortCode": 174379,
		"Password":         os.Getenv("Safaricom_Password"),
		"Timestamp":        "20231217153132",
		"TransactionType":  "CustomerPayBillOnline",
		"Amount":           payment.Cost,
		"PartyA":           25497408042,
		"PartyB":           174379,
		"PhoneNumber":      payment.CustomerPhone,
		"CallBackURL":      "https://medicare-t9y1.onrender.com/callback",
		"AccountReference": "Furahia Adventures",
		"TransactionDesc":  "Bike Ride",
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Return the payment details
	return &payment, nil
}



func handleCallback(c *fiber.Ctx) error {
	var callback struct {
		Body struct {
			StkCallback struct {
				MerchantRequestID string `json:"MerchantRequestID"`
				CheckoutRequestID string `json:"CheckoutRequestID"`
				ResultCode        int    `json:"ResultCode"`
				ResultDesc        string `json:"ResultDesc"`
				CallbackMetadata  struct {
					Item []struct {
						Name  string      `json:"Name"`
						Value interface{} `json:"Value"`
					} `json:"Item"`
				} `json:"CallbackMetadata"`
			} `json:"stkCallback"`
		} `json:"Body"`
	}

	// Parse the callback payload
	if err := c.BodyParser(&callback); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Extract relevant fields
	var (
		amount          float64
		transactionID   string
		phoneNumber     string
		transactionDate string
	)

	for _, item := range callback.Body.StkCallback.CallbackMetadata.Item {
		switch item.Name {
		case "Amount":
			amount = item.Value.(float64)
		case "MpesaReceiptNumber":
			transactionID = item.Value.(string)
		case "PhoneNumber":
			phoneNumber = fmt.Sprintf("%v", item.Value) // Convert to string
		case "TransactionDate":
			transactionDate = fmt.Sprintf("%v", item.Value) // Convert to string
		}
	}

	// Determine payment status
	paymentStatus := "Failed"
	if callback.Body.StkCallback.ResultCode == 0 {
		paymentStatus = "Completed"
	}

	// Create a new Payments struct
	payment := model.Payments{
		ID:              uuid.New(),
		Cost:            amount,
		PaymentMethod:   "M-Pesa",
		TransactionID:   transactionID,
		PaymentStatus:   paymentStatus,
		CustomerPhone:   phoneNumber,
		AccountReference: "Order#123", // You can dynamically set this based on your logic
		TransactionDesc: "Payment for Order #123",
		TransactionDate: transactionDate,
	}

	// Return the payment details
	return c.JSON(payment)
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