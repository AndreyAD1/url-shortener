package middlewares

import (
	"errors"
	"net/http"
)

const userCookieName = "url-shortener"

func getCookieValue() string {
	return "encrypted user ID"
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(userCookieName)
		if errors.Is(http.ErrNoCookie, err) {
			cookieValue := getCookieValue()
			newCookie := http.Cookie{
				Name:  userCookieName,
				Value: cookieValue,
			}
			r.AddCookie(&newCookie)
			http.SetCookie(w, &newCookie)
		}
		next.ServeHTTP(w, r)
	})
}
