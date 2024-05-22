package transfer

import "testbanc/pkg/account"

type TransferMoney struct {
	FromAccount *account.Account `json:"fromAccount"`
	ToAccount   *account.Account `json:"toAccount"`
	Amount      float64          `json:"amount"`
}
