package middlewares

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/iamsayantan/larasockets"
	"github.com/iamsayantan/larasockets/server/handlers/dto"
	"github.com/iamsayantan/larasockets/server/rendering"
	"net/http"
)

const keyAuthAppId = 0

func NewAuthMiddleware(am larasockets.ApplicationManager) *AuthMiddleware {
	return &AuthMiddleware{appManager: am}
}

type AuthMiddleware struct {
	appManager larasockets.ApplicationManager
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appId := chi.URLParam(r, "appId")
		app := am.appManager.FindById(appId)

		if app == nil {
			rendering.RenderError(w, "Unauthenticated Access", http.StatusUnauthorized)
			return
		}

		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			rendering.RenderError(w, "Unauthenticated Access", http.StatusUnauthorized)
			return
		}

		tokenClaims := &dto.ApplicationAuthorizationClaims{}
		token, err := jwt.ParseWithClaims(accessToken, tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Secret()), nil
		})

		if err != nil {
			rendering.RenderError(w, "Unauthenticated Access", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			rendering.RenderError(w, "Unauthenticated Access", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, keyAuthAppId, tokenClaims.AppId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthenticatedAppIdFromContext(ctx context.Context) string {
	return ctx.Value(keyAuthAppId).(string)
}
