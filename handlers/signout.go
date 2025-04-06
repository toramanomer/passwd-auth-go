package handlers

import (
	"net/http"

	"github.com/toramanomer/passwd-auth-go/core/repository"
	sessioncookie "github.com/toramanomer/passwd-auth-go/core/sessionCookie"
)

type SignoutController struct {
	UserManagementRepo repository.UserManagementRepository
}

func (c *SignoutController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST or OPTIONS
	// and set the appropriate CORS headers
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(sessioncookie.GetName())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c.UserManagementRepo.DeleteSession(r.Context(), cookie.Value)

	http.SetCookie(w, sessioncookie.Expire())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
