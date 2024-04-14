package middleware

import (
	"banner/internal/session"
	"net/http"
)

// Auth проверяет, авторизован ли пользователь
func Auth(snm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, rqt *http.Request) {
		sess, err := session.SessionFromContext(rqt.Context())
		if err != nil || sess == nil {
			if sess = snm.CheckToken(wrt, rqt); sess != nil {
				ctx := session.ContextWithSession(rqt.Context(), sess)
				next.ServeHTTP(wrt, rqt.WithContext(ctx))
				return
			}
		}
	})
}
