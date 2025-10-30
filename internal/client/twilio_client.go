package client

import (
	"fmt"
    "os"
    "log"
	"github.com/twilio/twilio-go"
	twilio_api "github.com/twilio/twilio-go/rest/api/v2010"
    "github.com/Bilal-Cplusoft/sunready/utils"
    "strings"
)

type TwilioClient struct {
	client        *twilio.RestClient
	fromNumber    string
}


func InitializeTwilio() *TwilioClient {
	accountSID := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_AUTH")
	fromNumber := os.Getenv("TWILIO_FROM")
	if accountSID == "" || authToken == "" || fromNumber == "" {
     log.Fatal("Twilio SID, Auth Token, or From Number not set in environment")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioClient{
		client:        client,
		fromNumber:    fromNumber,
	}
}

func (tc *TwilioClient) SendOTP(phoneNumber string) (string, error) {
	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	body := fmt.Sprintf("Your SunReady verification code is %s", otp)

	params := &twilio_api.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(tc.fromNumber)
	params.SetBody(body)

	_, err = tc.client.Api.CreateMessage(params)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP SMS: %w", err)
	}

	fmt.Printf("OTP sent to %s: %s\n", phoneNumber, otp)
	return otp, nil
}
