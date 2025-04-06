package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/toramanomer/passwd-auth-go/core/emailverification"
	"github.com/toramanomer/passwd-auth-go/core/mailer"
	"github.com/toramanomer/passwd-auth-go/core/model"
	"github.com/toramanomer/passwd-auth-go/core/password"
	"github.com/toramanomer/passwd-auth-go/core/repository"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupController struct {
	UserManagementRepo        repository.UserManagementRepository
	EmailVerificationStrategy emailverification.EmailVerificationStrategy
	Mailer                    mailer.Mailer
}

func (c *SignupController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST or OPTIONS
	// and set the appropriate CORS headers
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Parse body
	body, _ := io.ReadAll(r.Body)
	var signupRequest signupRequest
	if err := json.Unmarshal(body, &signupRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		email = strings.ToLower(strings.TrimSpace(signupRequest.Email))
		pass  = password.HashPassword(signupRequest.Password)
	)

	// Validate email and password
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if pass == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	} else if utf8.RuneCountInString(pass) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	} else if utf8.RuneCountInString(pass) > 128 {
		http.Error(w, "Password must be less than 128 characters long", http.StatusBadRequest)
		return
	}

	var (
		user                                 = model.NewUser(email, pass)
		rawCode, protectedCode, strategyName = c.EmailVerificationStrategy.GenerateCode()
		ev                                   = model.NewEmailVerification(user.ID, protectedCode, strategyName)
	)

	// Create user and email verification
	if err := c.UserManagementRepo.CreatePendingUser(r.Context(), user, ev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send email verification code
	if err := c.Mailer.SendVerificationEmail(email, rawCode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{ "success": true }`))
}
