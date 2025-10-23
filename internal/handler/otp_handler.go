package handler

import (
	"github.com/Bilal-Cplusoft/sunready/internal/client"
	"net/http"
	"sync"
	"time"
	"errors"
	"fmt"
	"encoding/json"
)

type OtpHandler struct {
	twilioClient *client.TwilioClient
	sendgridClient *client.SendGridClient
	otpStore      map[string]OTPData
	mu            sync.RWMutex
	otpExpiration time.Duration
}

type OTPData struct {
	Code      string
	ExpiresAt time.Time
	Attempts  int
}

type MessageResponse struct {
	Message string `json:"message"`
}


func NewOtpHandler(twilioClient *client.TwilioClient, sendgridClient *client.SendGridClient) *OtpHandler {
	return &OtpHandler{
		twilioClient: twilioClient,
		sendgridClient: sendgridClient,
		otpStore:      make(map[string]OTPData),
		otpExpiration: 10 * time.Minute,
	}
}


// SendOTP godoc
// @Summary      Send OTP to a phone number
// @Description  Sends a one-time password (OTP) via SMS to the specified phone number using Twilio.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        phone   query     string  true  "Phone number with country code (e.g. +923001234567)"
// @Param        email   query     string  true  "Email Address (e.g. test@example.com)"
// @Success 200 {string} string "OTP sent successfully"
// @Failure 400 {string} string "Missing or invalid phone parameter"
// @Failure 500 {string} string "Failed to send OTP"
// @Router       /api/otp/send [get]
func (h *OtpHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	phone := r.URL.Query().Get("phone")
	email := r.URL.Query().Get("email")

	if phone == "" || email == "" {
		http.Error(w, "missing phone or email parameter", http.StatusBadRequest)
		return
	}

	smsOTP, err := h.twilioClient.SendOTP(phone)
	if err != nil {
		http.Error(w, "failed to send phone OTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	emailOTP, err := h.sendgridClient.SendOTP(email)
	if err != nil {
		http.Error(w, "failed to send email OTP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	expiry := time.Now().Add(h.otpExpiration)
	h.otpStore[phone] = OTPData{Code: smsOTP, ExpiresAt: expiry, Attempts: 0}
	h.otpStore[email] = OTPData{Code: emailOTP, ExpiresAt: expiry, Attempts: 0}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTPs sent successfully",
	})
}


// VerifyOTP godoc
// @Summary      Verify phone and email OTP codes
// @Description  Verifies one-time passwords (OTPs) sent to both phone and email.
// @Tags         OTP
// @Accept       json
// @Produce      json
// @Param        phone       query     string  true  "Phone number with country code (e.g. +923001234567)"
// @Param        sms_otp     query     string  true  "OTP code received via SMS"
// @Param        email       query     string  true  "Email address used for OTP verification"
// @Param        email_otp   query     string  true  "OTP code received via Email"
// @Success 200 {object} MessageResponse "OTP verified successfully"
// @Failure 400 {object} MessageResponse "Missing or invalid parameters / OTP verification failed"
// @Router       /api/otp/verify [get]
func (h *OtpHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	phone := r.URL.Query().Get("phone")
	smsOTP := r.URL.Query().Get("sms_otp")
	email := r.URL.Query().Get("email")
	emailOTP := r.URL.Query().Get("email_otp")

	if phone == "" || smsOTP == "" || email == "" || emailOTP == "" {
		http.Error(w, "missing required query parameters", http.StatusBadRequest)
		return
	}

	if err := h.verifySingleOTP(phone, smsOTP); err != nil {
		http.Error(w, "phone OTP: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.verifySingleOTP(email, emailOTP); err != nil {
		http.Error(w, "email OTP: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "OTP verified successfully",
	})
}

func (h *OtpHandler) verifySingleOTP(key, otp string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, ok := h.otpStore[key]
	if !ok {
		return errors.New("no OTP found")
	}

	now := time.Now()
	switch {
	case now.After(data.ExpiresAt):
		delete(h.otpStore, key)
		return errors.New("OTP expired")
	case data.Attempts >= 3:
		delete(h.otpStore, key)
		return errors.New("maximum attempts exceeded")
	case data.Code != otp:
		data.Attempts++
		h.otpStore[key] = data
		return fmt.Errorf("invalid OTP (attempt %d/3)", data.Attempts)
	default:
		delete(h.otpStore, key)
		return nil
	}
}
