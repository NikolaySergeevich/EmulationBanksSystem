package repository

import (
	"context"
	"testbanc/pkg/account"
)

type Repository interface{
	AddGovernment(ctx context.Context, countryAccount *account.Account, destructionAccount *account.Account) error
	AddAccount(ctx context.Context, account *account.Account) error
	AddMoney(ctx context.Context, accountNumber string, count int) error
	DestroyMoney(ctx context.Context, accountNumber string, count int) error
	GetAllAccount(ctx context.Context) ([]account.GetAccount, error)
}