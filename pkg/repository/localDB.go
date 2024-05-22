package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testbanc/pkg/account"
)

var _ Repository = (*LocalDB)(nil) // LocalDB соответствует интерфейсу

func NewLocalDB() Repository {
	return &LocalDB{
		Accounts:           make(map[string]*account.Account),
		CountryAccount:     &account.Account{},
		DestructionAccount: &account.Account{},
	}
}

type LocalDB struct {
	Accounts           map[string]*account.Account
	Mu                 sync.Mutex
	CountryAccount     *account.Account
	DestructionAccount *account.Account
}

func (l *LocalDB) AddGovernment(ctx context.Context, countryAccount *account.Account, destructionAccount *account.Account) error{
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, ok := l.Accounts[countryAccount.AccountNumber]; !ok {
		l.Accounts[countryAccount.AccountNumber] = countryAccount
		l.CountryAccount = countryAccount
	}else {
		return errors.New("db: account with Country such a number already exists")
	}
	if _, ok := l.Accounts[destructionAccount.AccountNumber]; !ok {
		l.Accounts[destructionAccount.AccountNumber] = destructionAccount
		l.DestructionAccount = destructionAccount
	}else {
		return errors.New("db: account with Destruction such a number already exists")
	}
	return nil
}

func (l *LocalDB) AddAccount(ctx context.Context, account *account.Account) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, ok := l.Accounts[account.AccountNumber]; !ok {
		l.Accounts[account.AccountNumber] = account
		return nil
	}
	fmt.Println(len(l.Accounts))
	return errors.New("db: account with such a number already exists")
}

func (l *LocalDB) AddMoney(ctx context.Context, accountNumber string, count int) error {
	return nil

}

func (l *LocalDB) DestroyMoney(ctx context.Context, accountNumber string, count int) error {
	return nil

}

func (l *LocalDB) GetAllAccount(ctx context.Context) ([]account.GetAccount, error) {
	res := make([]account.GetAccount, 0, len(l.Accounts)+2)
	res = append(res, account.GetAccount{AccountNumber: l.CountryAccount.AccountNumber, Balance: l.CountryAccount.Balance, Active: l.CountryAccount.Active})
	res = append(res, account.GetAccount{AccountNumber: l.DestructionAccount.AccountNumber, Balance: l.DestructionAccount.Balance, Active: l.DestructionAccount.Active})

	for _, v := range l.Accounts {
		res = append(res, account.GetAccount{AccountNumber: v.AccountNumber, Balance: v.Balance, Active: v.Active})
	}

	if len(res) != len(l.Accounts)+2 {
		return nil, errors.New("db: get all account")
	}
	return res, nil

}
