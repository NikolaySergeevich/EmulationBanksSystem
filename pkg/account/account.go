package account

import (
	"sync"
	"testbanc/pkg/tools"
	"time"
)

// Банковский счет.
type Account struct {
	AccountNumber tools.Iban
	Mu            sync.Mutex
	Balance       float64
	Active        bool
	CreatedAt     time.Time 
	UpdatedAt     time.Time 
}


type GetAccount struct {
	AccountNumber string
	Balance       float64
	Active        bool
}
