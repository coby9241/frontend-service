package transactions

import (
	"time"

	"github.com/coby9241/frontend-service/internal/models/users"
	"github.com/jinzhu/gorm"
)

// Transaction holds all necessary information for the transaction
type Transaction struct {
	gorm.Model
	UserID               *uint
	User                 users.User
	TransactionAmount    float64
	TrackingNumber       string
	CancelledReason      string
	PaymentProviderRefID string
	PaymentMethod        string
	ConfirmedAt          *time.Time
	CompletedAt          *time.Time
	CancelledAt          *time.Time
}
