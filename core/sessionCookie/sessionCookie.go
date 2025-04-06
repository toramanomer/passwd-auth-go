package sessioncookie

import "net/http"

const sessionCookieName = "__Host-s"
const maxAge = 60 * 10 // 10 minutes in seconds

func Create(sessionID string) *http.Cookie {

	return &http.Cookie{
		Name:   sessionCookieName,
		Value:  sessionID,
		Quoted: false,

		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func Expire() *http.Cookie {
	cookie := Create("")
	cookie.MaxAge = -1

	return cookie
}
