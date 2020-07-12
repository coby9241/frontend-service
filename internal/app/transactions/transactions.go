package transactions

import (
	"strings"

	"github.com/coby9241/frontend-service/internal/models/transactions"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

// New instantiates a new App for transaction management
func New() *App {
	return &App{}
}

// App home app
type App struct {
}

// ConfigureAdmin configures the admin page for CRUD of transactions
func (App) ConfigureAdmin(adm *admin.Admin) {
	adm.AddMenu(&admin.Menu{Name: "Transaction Management", Priority: 2})

	// Add Transaction
	transaction := adm.AddResource(&transactions.Transaction{}, &admin.Config{Menu: []string{"Transaction Management"}})
	transaction.Meta(&admin.Meta{Name: "ConfirmedAt", Type: "date"})

	// Scopes for Transaction
	for _, state := range []string{"pending", "validated", "cancelled"} {
		var state = state
		transaction.Scope(&admin.Scope{
			Name:  state,
			Label: strings.Title(strings.Replace(state, "_", " ", -1)),
			Group: "Transaction Status",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(transactions.Transaction{State: state})
			},
		})
	}
}
