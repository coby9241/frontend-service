package encryptor

import (
	"sync"

	"github.com/coby9241/frontend-service/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	instance *BcryptEncryptor
	once     sync.Once
)

// BcryptEncryptor is
type BcryptEncryptor struct {
	Cost int
}

// GetInstance returns a BcryptEncryptor pointer to retrieve environment variables
func GetInstance() *BcryptEncryptor {
	once.Do(func() {
		bcryptCost := bcrypt.DefaultCost
		if config.GetInstance().BcryptCost != 0 {
			bcryptCost = config.GetInstance().BcryptCost
		}

		instance = &BcryptEncryptor{
			Cost: bcryptCost,
		}
	})

	return instance
}

// Digest is
func (be BcryptEncryptor) Digest(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), be.Cost)
	return string(hashedPassword), err
}

// Compare is
func (be BcryptEncryptor) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
