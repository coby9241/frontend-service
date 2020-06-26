package users

import (
	"fmt"
	"time"

	"github.com/coby9241/frontend-service/internal/config"
	"github.com/coby9241/frontend-service/internal/encryptor"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// User is
type User struct {
	gorm.Model

	UID               string `gorm:"column:uid" sql:"type:varchar;unique_index:uix_identity_uid;index:idx_identity_uid_provider" valid:"required"` // username/email
	Provider          string `sql:"type:varchar;index:idx_identity_uid_provider" valid:"required,in(email)"`                                       // phone, email, github...
	PasswordHash      string `gorm:"column:password;not null"`
	Role              Role   `gorm:"foreignkey:ID"`
	UserID            string `sql:"type:varchar"` // user's name
	PasswordChangedAt *time.Time
}

// DisplayName shows the name of the user and if not found, the registered email in the admin dashboard
func (u User) DisplayName() string {
	if u.UserID == "" {
		return u.UID
	}

	return u.UserID
}

// BeforeCreate is
func (u User) BeforeCreate(tx *gorm.DB) (err error) {
	u.PasswordHash, err = encryptor.GetInstance().Digest(u.PasswordHash)
	timeNow := time.Now()
	u.PasswordChangedAt = &timeNow

	return
}

// ComparePassword is
func (u *User) ComparePassword(password string) error {
	return encryptor.GetInstance().Compare(u.PasswordHash, password)
}

// IssueJwtTokenSet is
func (u *User) IssueJwtTokenSet(jwtKey interface{}) (*TokenSet, error) {
	timeNow := time.Now()
	expiresAt := timeNow.AddDate(1, 0, 0).Unix() // one year

	claims := &Claims{
		PasswordChangedAt: u.PasswordChangedAt.Unix(),
		StandardClaims: jwt.StandardClaims{
			Subject:   u.UID,
			Issuer:    fmt.Sprintf("frontend-service-%s", config.GetInstance().AppEnv),
			IssuedAt:  timeNow.Unix(),
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &TokenSet{
		Token:         tokenString,
		ExpiresAtUnix: expiresAt,
	}, nil
}
