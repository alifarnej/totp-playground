package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

const issuer = "MyCompany"

var secretKey string

// Handler to generate QR code with accountName as query parameter
func generateQRCode(w http.ResponseWriter, r *http.Request) {
	accountName := r.URL.Query().Get("accountName")
	if accountName == "" {
		http.Error(w, "missing accountName parameter", http.StatusBadRequest)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		http.Error(w, "failed to generate TOTP key", http.StatusInternalServerError)
		return
	}

	secretKey = key.Secret() // Store the key for validation
	otpURL := key.URL()
	png, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(png)
}

// Handler to validate the TOTP code
func validateTOTP(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "missing code parameter", http.StatusBadRequest)
		return
	}

	valid := totp.Validate(code, secretKey)
	if !valid {
		http.Error(w, "invalid code", http.StatusUnauthorized)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("the code is valid"))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/qr", generateQRCode)
	r.Get("/v/{code}", validateTOTP)

	fmt.Println("starting server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
