package authmiddleware

import (
	"context"
	"net/http"
	"strings"

	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/shared"
)

type contextKey string

const userIDKey contextKey = "userID"

func UserID(ctx context.Context) (int32, bool) {
	id, ok := ctx.Value(userIDKey).(int32)
	return id, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			errormiddleware.HandleError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFromRequest(r *http.Request) (int32, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return 0, apierr.NewUnauthorized("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return 0, apierr.NewUnauthorized("invalid authorization header")
	}

	return shared.ParseToken(parts[1])
}
