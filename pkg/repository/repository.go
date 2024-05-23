package repository

import (
	"context"
	"testbanc/pkg/account"
	"testbanc/pkg/tools"
)

type Repository interface{
	AddGovernment(ctx context.Context, countryAccount *account.Account, destructionAccount *account.Account) error //+
	AddAccount(ctx context.Context, account *account.Account) error//+
	Emition(ctx context.Context, count float64) error//+
	DestroyMoney(ctx context.Context, accountNumber string, count float64) error//+
	TransferMoney(ctx context.Context, fromAccount, toAccount *account.Account, amount float64) error// +
	TransferMoneyJSON(ctx context.Context, transfer []byte) error
	FindAccountIBAN(ctx context.Context, iban string) (*account.Account, error)
	GetAllAccount(ctx context.Context) ([]*account.GetAccount, error)//+
	GetCountryIBAN(ctx context.Context) (tools.Iban, error)//+
	GetDestroyIBAN(ctx context.Context) (tools.Iban, error)//+
}