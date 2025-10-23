package client

import (
	"fmt"
	"log"
	"os"
    "github.com/Bilal-Cplusoft/sunready/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	client *sendgrid.Client
	from   *mail.Email
}

func InitializeSendGrid() *SendGridClient {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		log.Fatal("SendGrid API key or FROM email not set in environment")
	}

	client := sendgrid.NewSendClient(apiKey)
	from := mail.NewEmail("SunReady Team", fromEmail)

	return &SendGridClient{
		client: client,
		from:   from,
	}
}


func (sg *SendGridClient) SendWelcomeEmail(toEmail, name string) error {
	to := mail.NewEmail(name, toEmail)
	subject := "Welcome to SunReady!"
	plainTextContent := fmt.Sprintf("Hello %s,\n\nWelcome to SunReady! We're excited to have you on board.", name)
	htmlContent := fmt.Sprintf("<strong>Hello %s,</strong><br><br>Welcome to SunReady! We're excited to have you on board.", name)

	message := mail.NewSingleEmail(sg.from, subject, to, plainTextContent, htmlContent)
	response, err := sg.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	fmt.Printf("Welcome email sent to %s successfully\n", toEmail)
	return nil
}

func (sg *SendGridClient) SendOTP(toEmail string) (string,error) {
	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}
	to := mail.NewEmail("", toEmail)
	subject := "Your SunReady Verification Code"

	plainTextContent := fmt.Sprintf(
		"Your verification code is: %s\n\nThis code will expire in 10 minutes.\n\nIf you didn't request this code, please ignore this email.",
		otp,
	)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Verification Code</h2>
			<p>Your verification code is:</p>
			<div style="background-color: #f4f4f4; padding: 15px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; margin: 20px 0;">
				%s
			</div>
			<p style="color: #666;">This code will expire in 10 minutes.</p>
			<p style="color: #999; font-size: 12px;">If you didn't request this code, please ignore this email.</p>
		</div>
	`, otp)

	message := mail.NewSingleEmail(sg.from, subject, to, plainTextContent, htmlContent)

	response, err := sg.client.Send(message)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("failed to send OTP email, status code: %d, body: %s", response.StatusCode, response.Body)
	}
	fmt.Printf("OTP email sent to %s: %v \n", toEmail,otp)
	return otp, nil
}
