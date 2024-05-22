package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"testbanc/pkg/account"
	"testbanc/pkg/tools"
	"testbanc/pkg/transfer"
	"time"
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

// Добавление счёта государства и для уничтожения в условную БД
func (l *LocalDB) AddGovernment(ctx context.Context, countryAccount *account.Account, destructionAccount *account.Account) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, ok := l.Accounts[countryAccount.AccountNumber.Iban]; !ok {
		l.Accounts[countryAccount.AccountNumber.Iban] = countryAccount
		l.CountryAccount = countryAccount
	} else {
		return errors.New("db: account with Country such a number already exists")
	}
	if _, ok := l.Accounts[destructionAccount.AccountNumber.Iban]; !ok {
		l.Accounts[destructionAccount.AccountNumber.Iban] = destructionAccount
		l.DestructionAccount = destructionAccount
	} else {
		return errors.New("db: account with Destruction such a number already exists")
	}
	return nil
}

// Добавление счёта в систему хранения платёжной системы
func (l *LocalDB) AddAccount(ctx context.Context, account *account.Account) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, ok := l.Accounts[account.AccountNumber.Iban]; !ok {
		l.Accounts[account.AccountNumber.Iban] = account
		return nil
	}
	fmt.Println(len(l.Accounts))
	return errors.New("db: account with such a number already exists")
}

// Вернёт структуру Iban. А именно номер счёта, под которым зарегистрировано государство
func (l *LocalDB) GetCountryIBAN(ctx context.Context) (tools.Iban, error) {
	if l.CountryAccount == nil {
		return tools.Iban{}, errors.New("empty country account")
	}
	if _, ok := interface{}(l.CountryAccount.AccountNumber).(tools.Iban); ok {
		return l.CountryAccount.AccountNumber, nil
	}
	return tools.Iban{}, errors.New("the IBAN is of a different type")
}

// Вернёт структуру Iban. А именно номер счёта, под которым зарегистрирован счёт для уничтожения денег
func (l *LocalDB) GetDestroyIBAN(ctx context.Context) (tools.Iban, error) {
	if l.DestructionAccount == nil {
		return tools.Iban{}, errors.New("empty destruction account")
	}
	if _, ok := interface{}(l.DestructionAccount.AccountNumber).(tools.Iban); ok {
		return l.DestructionAccount.AccountNumber, nil
	}
	return tools.Iban{}, errors.New("the IBAN is of a different type")
}

// Выполняет эмиссию с государственным счётом
func (l *LocalDB) Emition(ctx context.Context, count float64) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			log.Println("the issue was not carried out")
			l.Mu.Unlock()
		}
	}()

	l.CountryAccount.Balance += count
	l.CountryAccount.UpdatedAt = time.Now()
	return nil
}

// Уничтожение денег с указанного счёта.
func (l *LocalDB) DestroyMoney(ctx context.Context, accountNumber string, count float64) error {
	if accountNumber == l.DestructionAccount.AccountNumber.Iban {
		return errors.New("this operation is not possible")
	}
	acc, ok := l.Accounts[accountNumber]
	if !ok {
		return errors.New("there is no account with this number")
	}
	if acc.Balance < count {
		return fmt.Errorf("insufficient funds. Maximum amount to destroy: %.2f", acc.Balance)
	}
	l.Mu.Lock()
	defer l.Mu.Unlock()
	acc.Balance -= count
	acc.UpdatedAt = time.Now()
	l.DestructionAccount.Balance += count
	l.DestructionAccount.UpdatedAt = time.Now()
	return nil
}

// осуществлять перевод заданной суммы денег между двумя указанными счетами
func (l *LocalDB) TransferMoney(ctx context.Context, fromAccount, toAccount *account.Account, amount float64) error {
	err := l.CheckTransfer(fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	l.Mu.Lock()
	defer l.Mu.Unlock()
	fromAccount.Balance -= amount
	fromAccount.UpdatedAt = time.Now()
	toAccount.Balance += amount
	toAccount.UpdatedAt = time.Now()
	return nil
}

// осуществлять перевод заданной суммы денег между двумя указанными счетами. JSON версия
func (l *LocalDB) TransferMoneyJSON(ctx context.Context, transf []byte) error {
	var tr transfer.TransferMoney
	if err := json.Unmarshal(transf, &tr); err != nil {
		return fmt.Errorf("unmarshal transfer: %w", err)
	}
	err := l.CheckTransfer(tr.FromAccount, tr.ToAccount, tr.Amount)
	if err != nil {
		return err
	}
	l.Mu.Lock()
	defer l.Mu.Unlock()
	tr.FromAccount.Balance -= tr.Amount
	tr.FromAccount.UpdatedAt = time.Now()
	tr.ToAccount.Balance += tr.Amount
	tr.ToAccount.UpdatedAt = time.Now()
	// Блок ниже нужен, если обменивающимися счетами являются гос. счета
	switch tr.FromAccount.AccountNumber.Iban {
	case l.CountryAccount.AccountNumber.Iban:
		l.CountryAccount = tr.FromAccount
	case l.DestructionAccount.AccountNumber.Iban:
		l.DestructionAccount = tr.FromAccount
	}
	switch tr.ToAccount.AccountNumber.Iban{
	case l.CountryAccount.AccountNumber.Iban:
		l.CountryAccount = tr.ToAccount
	case l.DestructionAccount.AccountNumber.Iban:
		l.DestructionAccount = tr.ToAccount
	}
	return nil
}

// Проверяет валидность трансфера между двумя счетами
func (l *LocalDB) CheckTransfer(fromAccount, toAccount *account.Account, amount float64) error {
	if fromAccount.AccountNumber.Iban == l.DestructionAccount.AccountNumber.Iban {
		return errors.New("money cannot be returned from the account for destruction")
	}
	// fromAcc, ok := l.Accounts[fromAccount]
	// if !ok {
	// 	return &account.Account{}, &account.Account{}, errors.New("there is no account with fromAccount number")
	// }
	// toAcc, ok := l.Accounts[toAccount]
	// if !ok {
	// 	return &account.Account{}, &account.Account{}, errors.New("there is no account with toAccount number")
	// }
	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds. Maximum amount to transfer: %.2f", fromAccount.Balance)
	}
	return nil
}

// Ищет счёт по его номеру IBAN
func (l *LocalDB) FindAccountIBAN(ctx context.Context, iban string) (*account.Account, error) {
	Acc, ok := l.Accounts[iban]
	if !ok {
		return &account.Account{}, errors.New("there is no account with this number")
	}
	return Acc, nil
}

// Вернёт список со всеми имеющимися счетами в системе, включая служебные. Номер, остато на счёте, статус счёта.
func (l *LocalDB) GetAllAccount(ctx context.Context) ([]account.GetAccount, error) {
	res := make([]account.GetAccount, 0, len(l.Accounts)+2)
	res = append(res, account.GetAccount{AccountNumber: l.CountryAccount.AccountNumber.Iban, Balance: l.CountryAccount.Balance, Active: l.CountryAccount.Active})
	res = append(res, account.GetAccount{AccountNumber: l.DestructionAccount.AccountNumber.Iban, Balance: l.DestructionAccount.Balance, Active: l.DestructionAccount.Active})

	for _, v := range l.Accounts {
		if v.AccountNumber.Iban == l.CountryAccount.AccountNumber.Iban || v.AccountNumber.Iban == l.DestructionAccount.AccountNumber.Iban {
			continue
		}
		res = append(res, account.GetAccount{AccountNumber: v.AccountNumber.Iban, Balance: v.Balance, Active: v.Active})
	}

	if len(res) != len(l.Accounts) {
		return nil, errors.New("db: get all account")
	}
	return res, nil

}
