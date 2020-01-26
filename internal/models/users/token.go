package users

import (
	"errors"
	"frontend-project/internal/config"

	"github.com/dgrijalva/jwt-go"
)

type (
	// Claims is
	Claims struct {
		PasswordChangedAt int64 `json:"xpca,omitempty"`
		jwt.StandardClaims
	}

	// TokenSet is
	TokenSet struct {
		Token         string `json:"token"`
		ExpiresAtUnix int64  `json:"expires_at_unix"`
	}
)

// GetClaims is
func GetClaims(jwtToken string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetInstance().JwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
