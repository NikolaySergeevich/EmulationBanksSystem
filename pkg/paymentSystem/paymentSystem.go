package paymentsystem

import (
	"context"
	"encoding/json"
	"fmt"
	"testbanc/pkg/account"
	"testbanc/pkg/repository"
	"testbanc/pkg/tools"
	"time"
)

// Платёжная система.
type PaymentSystem struct {
	accounts           repository.Repository
}

// Создание новой платёжной системы.
func NewPaymentSystem(ctx context.Context, db repository.Repository) (PaymentSystem, error) {
	countryAccount := account.Account{
		AccountNumber: tools.GenerateIBAN(),
		Balance:       0,
		Active:        true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	destructionAccount := account.Account{
		AccountNumber: tools.GenerateIBAN(),
		Balance:       0,
		Active:        true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.AddGovernment(ctx, &countryAccount, &destructionAccount); err != nil{
		return PaymentSystem{}, fmt.Errorf("new payment: %w", err)
	}

	return PaymentSystem{
		accounts:           db,
	}, nil
}

// func (ps *PaymentSystem) GetCountryNumber()

// Создание нового платёжного счёта.
func (ps *PaymentSystem) CreateAccount(ctx context.Context) error {
	acc := account.Account{
		AccountNumber: tools.GenerateIBAN(),
		Balance:       0,
		Active:        true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := ps.accounts.AddAccount(ctx, &acc); err != nil{
		return fmt.Errorf("create acc: %w", err)
	}
	return nil
}

// Добавляет указанную сумму на государственный счет.
// func (ps *PaymentSystem) EmitMoney(amount float64) {
// 	ps.countryAccount.Mu.Lock()
// 	defer ps.countryAccount.Mu.Unlock()
// 	ps.countryAccount.Balance += amount
// }

// // Списывает указанную сумму с государственного счета.
// func (ps *PaymentSystem) DestroyMoney(amount float64) {
// 	ps.countryAccount.Mu.Lock()
// 	defer ps.countryAccount.Mu.Unlock()
// 	ps.countryAccount.Balance -= amount
// }

// Переводит указанную сумму с одного счета на другой.
// func (ps *PaymentSystem) TransferMoney(fromAccount, toAccount string, amount float64) {
// 	if from, ok1 := ps.accounts[fromAccount]; ok1 {
// 		if to, ok2 := ps.accounts[toAccount]; ok2 {
// 			from.Mu.Lock()
// 			to.Mu.Lock()
// 			defer from.Mu.Unlock()
// 			defer to.Mu.Unlock()
// 			from.Balance -= amount
// 			to.Balance += amount
// 		}
// 	}
// }

// Возвращает список всех учетных записей.
func (ps *PaymentSystem) GetAccounts(ctx context.Context) (accounts string, err error) {
	acc, err := ps.accounts.GetAllAccount(ctx)
	if err != nil{
		return "", err
	}

	jsonData, err := json.Marshal(acc)
	if err != nil{
		return accounts, fmt.Errorf("get all accounts: %w", err)
	}
	return string(jsonData), nil
}
