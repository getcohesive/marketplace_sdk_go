package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/getcohesive/marketplace_sdk_go/pkg/common/utils"
)

type AuthDetails struct {
	UserID                 *int
	UserName               *string
	UserEmail              *string
	Role                   *string
	WorkspaceID            *int
	WorkspaceName          *string
	InstanceID             *int
	IsInTrial              *bool
	CurrentPeriodStartedAt *string
	CurrentPeriodEndsAt    *string
	TrialItemsCount        *int
	ItemsPerUnit           *int
}

func (authDetails AuthDetails) IsUserLoggedIn() bool {
	return authDetails.UserID != nil
}

func ValidateToken(tokenString string, appSecret string) (*AuthDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, AuthError("unexpected signing method. need HMAC")
		}
		return []byte(appSecret), nil
	})
	if err != nil {
		return nil, AuthError(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		authDetails := &AuthDetails{}
		authDetails.Role = utils.String(claims["role"])
		authDetails.WorkspaceID = utils.Int(claims["workspace_id"])
		authDetails.WorkspaceName = utils.String(claims["workspace_name"])
		authDetails.InstanceID = utils.Int(claims["instance_id"])
		authDetails.IsInTrial = utils.Bool(claims["is_in_trial"])
		authDetails.UserID = utils.Int(claims["user_id"])
		authDetails.UserEmail = utils.String(claims["user_name"])
		authDetails.UserEmail = utils.String(claims["user_email"])
		authDetails.CurrentPeriodStartedAt = utils.String(claims["current_period_started_at"])
		authDetails.CurrentPeriodEndsAt = utils.String(claims["current_period_ends_at"])
		authDetails.TrialItemsCount = utils.Int(claims["trial_items_count"])
		authDetails.ItemsPerUnit = utils.Int(claims["items_per_unit"])

		return authDetails, nil
	} else {
		return nil, AuthError("invalid token claims: " + token.Claims.Valid().Error())
	}
}
