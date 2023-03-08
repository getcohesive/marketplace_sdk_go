package middleware

import (
	"context"
	"net/http"
	"strings"

	cohesiveMarketplaceSDK "github.com/getcohesive/marketplace_sdk_go"
	"github.com/getcohesive/marketplace_sdk_go/pkg/authentication"
)

type AuthMiddleware struct {
	sdkClient cohesiveMarketplaceSDK.Client
}

func (mw *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, authentication.AuthError("auth header empty").Error(), http.StatusUnauthorized)
		return
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

	authDetails, err := mw.sdkClient.ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), authentication.AuthDetails{}, authDetails))
}

func NewAuthMiddleware(client cohesiveMarketplaceSDK.Client) *AuthMiddleware {
	return &AuthMiddleware{sdkClient: client}
}
