package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/toramanomer/passwd-auth-go/core/emailverification"
	"github.com/toramanomer/passwd-auth-go/core/mailer"
	"github.com/toramanomer/passwd-auth-go/core/model"
	"github.com/toramanomer/passwd-auth-go/core/repository"
)

type ResendVerificationController struct {
	UserManagementRepo        repository.UserManagementRepository
	EmailVerificationStrategy emailverification.EmailVerificationStrategy
	Mailer                    mailer.Mailer
}

type resendVerificationRequest struct {
	Email string `json:"email"`
}

func (c *ResendVerificationController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var req resendVerificationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		email = strings.ToLower(strings.TrimSpace(req.Email))
	)

	if email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := c.UserManagementRepo.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User could not be found", http.StatusBadRequest)
		return
	}

	if user.EmailVerifiedAt != nil {
		http.Error(w, "Already verified", http.StatusBadRequest)
		return
	}

	var (
		rawCode, protectedCode, strategyName = c.EmailVerificationStrategy.GenerateCode()
		ev                                   = model.NewEmailVerification(user.ID, protectedCode, strategyName)
	)

	if err := c.UserManagementRepo.CreateEmailVerification(r.Context(), ev); err != nil {
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
