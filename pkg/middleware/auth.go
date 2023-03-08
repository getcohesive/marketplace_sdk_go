package middleware

import (
	"context"
	"net/http"
	"strings"

	cohesiveMarketplaceSDK "github.com/getcohesive/marketplace_sdk_go"
	"github.com/getcohesive/marketplace_sdk_go/pkg/authentication"
)

type AuthMiddleware struct {
	sdkClient    cohesiveMarketplaceSDK.Client
	blockRequest bool
}

func (authMiddleware *AuthMiddleware) ParseAuthHeader(r *http.Request) (*authentication.AuthDetails, error) {
	if r.Method == http.MethodOptions {
		return nil, nil
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, authentication.AuthError("auth header empty")
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

	authDetails, err := authMiddleware.sdkClient.ValidateToken(token)
	if err != nil {
		return nil, authentication.AuthError("auth header validation failed" + err.Error())
	}
	return authDetails, nil
}

func (authMiddleware *AuthMiddleware) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		authDetails, err := authMiddleware.ParseAuthHeader(r)
		if err != nil && authMiddleware.blockRequest {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), authentication.AuthDetails{}, authDetails))
	}
}

func (authMiddleware *AuthMiddleware) HandlerFuncWithNext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			return
		}

		authDetails, err := authMiddleware.ParseAuthHeader(r)
		if err != nil && authMiddleware.blockRequest {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), authentication.AuthDetails{}, authDetails))
		next.ServeHTTP(w, r)
	}
}

func NewAuthMiddleware(client cohesiveMarketplaceSDK.Client, blockRequest bool) *AuthMiddleware {
	return &AuthMiddleware{sdkClient: client, blockRequest: blockRequest}
}
