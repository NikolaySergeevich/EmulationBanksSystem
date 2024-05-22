package account

import (
	"sync"
	"time"
)

// Банковский счет.
type Account struct {
	AccountNumber string
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
