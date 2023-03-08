package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type AuthDetails struct {
	UserID                  int
	UserName                string
	UserEmail               string
	Role                    string
	WorkspaceID             int
	WorkspaceName           string
	InstanceID              int
	CurrentPeriodStartedAt  string
	CurrentPeriodEndsAt     string
	IsInTrial               bool
	TrialItemsCount         int
}

func ValidateToken(tokenString string, appSecret string) (*AuthDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, AuthenticationError("unexpected signing method. need HMAC")
		}
		return []byte(appSecret), nil
	})
	if err != nil {
		return nil, AuthenticationError(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &AuthDetails{
			UserID:                  claims["user_id"].(int),
			UserName:                claims["user_name"].(string),
			UserEmail:               claims["user_email"].(string),
			Role:                    claims["role"].(string),
			WorkspaceID:             (claims["workspace_id"].(int)),
			WorkspaceName:           claims["workspace_name"].(string),
			InstanceID:              (claims["instance_id"].(int)),
			CurrentPeriodStartedAt:  getStringFromInterface(claims["current_period_started_at"]),
			CurrentPeriodEndsAt:     getStringFromInterface(claims["current_period_ends_at"]),
			IsInTrial:               claims["is_in_trial"].(bool),
			TrialItemsCount:        claims["trial_items_count"].(int),
		}, nil
	} else {
		return nil, AuthenticationError("invalid token claims: "+token.Claims.Valid().Error())
	}
}

func getStringFromInterface(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func getIntFromInterface(val interface{}) int {
	if val == nil {
		return 0
	}
	return int(val.(float64))
}
