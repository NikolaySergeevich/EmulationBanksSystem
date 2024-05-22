package paymentsystem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testbanc/pkg/account"
	"testbanc/pkg/repository"
	"testbanc/pkg/tools"
	"testbanc/pkg/transfer"
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
// Получение номера счета страны для "эмиссии".
func (ps *PaymentSystem) GetCountryIBAN(ctx context.Context) (string, error){
	iba, err := ps.accounts.GetCountryIBAN(ctx)
	if err != nil{
		return "", err
	}
	return iba.Iban, nil
}  

// Получение номера счета для "уничтожения" денег
func (ps *PaymentSystem) GetDestructionIBAN(ctx context.Context) (string, error){
	iba, err := ps.accounts.GetDestroyIBAN(ctx)
	if err != nil{
		return "", err
	}
	return iba.Iban, nil
}

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
func (ps *PaymentSystem) EmitMoney(ctx context.Context, amount float64) error {
	if err := ps.accounts.Emition(ctx, amount); err != nil{
		return err
	}
	return nil
}

// Списывает указанную сумму с указанного счета для уничтожения.
func (ps *PaymentSystem) DestroyMoney(ctx context.Context, accountNumber string, amount float64) error{
	if err := ps.accounts.DestroyMoney(ctx,accountNumber, amount); err != nil{
		return err
	}
	return nil
}

// Переводит указанную сумму с одного счета на другой.
func (ps *PaymentSystem) TransferMoney(ctx context.Context, fromAccount, toAccount string, amount float64) error{
	fromAcco, err := ps.accounts.FindAccountIBAN(ctx, fromAccount)
	if err != nil{
		return errors.New("there is no account with fromAccount number")
	}
	toAcco, err := ps.accounts.FindAccountIBAN(ctx, toAccount)
	if err != nil{
		return errors.New("there is no account with toAccount number")
	}
	if err := ps.accounts.TransferMoney(ctx, fromAcco, toAcco, amount); err != nil{
		return err
	}
	return nil
}

// Переводит указанную сумму с одного счета на другой. JSON версия
func (ps *PaymentSystem) TransferMoneyJSON(ctx context.Context, fromAccount, toAccount string, amount float64) error{
	var transf transfer.TransferMoney

	fromAc, err := ps.accounts.FindAccountIBAN(ctx, fromAccount)
	if err != nil{
		return errors.New("there is no account with fromAccount number")
	}
	toAc, err := ps.accounts.FindAccountIBAN(ctx, toAccount)
	if err != nil{
		return errors.New("there is no account with toAccount number")
	}
	
	transf.FromAccount = fromAc
	transf.ToAccount = toAc
	transf.Amount = amount

	data, err := json.Marshal(transf)
	if err != nil {
		return fmt.Errorf("marshal transfer: %w", err)
	}
	if err := ps.accounts.TransferMoneyJSON(ctx, data); err != nil{
		return err
	}
	return nil
}

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
