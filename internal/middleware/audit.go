package middleware

import (
	"net/http"
	//"os/user"
	"strconv"

	"github.com/chaitu35/claimeasy/internal/audit"
)

func InjectAuditUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from the request context or headers
		userId:=extractUserID(r)
		ctx:=audit.WithUserID(r.Context(), userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) uint {
	userIdStr:= r.Header.Get("X-User-ID")
	if userIdStr == "" {return 0}
	id, err := strconv.Atoi(userIdStr)
	if err != nil|| id < 0 {
		return 0
	}
	return uint(id)
}