package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

type contextKey string

const (
	AdminIDKey  contextKey = "admin_id"
	UsernameKey contextKey = "username"
)

var store *sessions.CookieStore

func InitAuth(sessionStore *sessions.CookieStore) {
	store = sessionStore
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "admin-session")
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		adminID, ok := session.Values["admin_id"]
		if !ok || adminID == nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		username, _ := session.Values["username"].(string)

		// Добавить данные в контекст
		ctx := context.WithValue(r.Context(), AdminIDKey, adminID)
		ctx = context.WithValue(ctx, UsernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAdminID(ctx context.Context) int {
	if id, ok := ctx.Value(AdminIDKey).(int); ok {
		return id
	}
	return 0
}

func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}
