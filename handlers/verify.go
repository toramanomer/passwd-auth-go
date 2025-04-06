package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/toramanomer/passwd-auth-go/core/emailverification"
	"github.com/toramanomer/passwd-auth-go/core/model"
	"github.com/toramanomer/passwd-auth-go/core/repository"
	sessioncookie "github.com/toramanomer/passwd-auth-go/core/sessionCookie"
)

type verifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerifyController struct {
	UserManagementRepo repository.UserManagementRepository
}

func (c *VerifyController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Parse body
	body, _ := io.ReadAll(r.Body)
	var verifyRequest verifyRequest
	if err := json.Unmarshal(body, &verifyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		email = strings.ToLower(strings.TrimSpace(verifyRequest.Email))
		code  = verifyRequest.Code
	)

	// Validate email and code
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	} else if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	user, ev, err := c.UserManagementRepo.GetUserAndEmailVerification(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.EmailVerifiedAt != nil {
		http.Error(w, "Already verified", http.StatusBadRequest)
		return
	}

	if ev.AttemptCount >= 3 {
		http.Error(w, "Attempt count exceeded, request new code.", http.StatusBadRequest)
		return
	} else if time.Now().UTC().After(ev.ExpiresAt) {
		http.Error(w, "The code has been expired, request new code.", http.StatusBadRequest)
		return
	}

	if !emailverification.VerifyCode(code, ev.VerificationCode) {
		err := c.UserManagementRepo.IncrementAttemptCount(r.Context(), ev.ID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Error(w, "The code does not match", http.StatusBadRequest)
		return
	}

	userSession := model.NewUserSession(user.ID, r.RemoteAddr, r.UserAgent())
	if err := c.UserManagementRepo.VerifyUserEmail(r.Context(), user.ID, ev.ID, userSession); err != nil {
		fmt.Printf("verify VerifyUserEmail %+v\n", err)
		http.Error(w, "Failed to verify", http.StatusInternalServerError)
		return
	}

	cookie := sessioncookie.Create(userSession.ID)
	http.SetCookie(w, cookie)
	w.Write([]byte(`{"success": true}`))
}
