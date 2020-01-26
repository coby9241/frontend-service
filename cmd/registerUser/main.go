package main

import (
	"flag"
	"time"

	"github.com/coby9241/frontend-service/internal/db"
	"github.com/coby9241/frontend-service/internal/encryptor"
	"github.com/coby9241/frontend-service/internal/models/users"
)

var (
	email    string
	password string
)

func init() {
	flag.StringVar(&email, "email", "", "email to register user")
	flag.StringVar(&password, "password", "", "password to register user")

	flag.Parse()
}

func main() {
	registerUser(email, password)
}

// RegisterUser registers a new user
func registerUser(email, password string) (*users.User, error) {
	passwordHash, err := encryptor.GetInstance().Digest(password)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()
	identity := &users.User{
		Provider:          "email",
		UID:               email,
		PasswordHash:      passwordHash,
		PasswordChangedAt: &timeNow,
	}

	if err = db.GetInstance().Create(identity).Error; err != nil {
		return nil, err
	}

	return identity, nil
}
