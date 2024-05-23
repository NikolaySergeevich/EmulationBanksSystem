package repository_test

import (
	"context"
	"testbanc/pkg/account"
	"testbanc/pkg/repository"
	"testbanc/pkg/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGovernment(t *testing.T) {
	countryAcc := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
	}
	destructionAcc := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
	}
	testCases := []struct {
		desc               string
		db                 *repository.LocalDB
		countryAccount     *account.Account
		destructionAccount *account.Account
		wait               bool
	}{
		{
			desc: "Success",
			db: &repository.LocalDB{
				Accounts: make(map[string]*account.Account),
			},
			countryAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
			},
			destructionAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY98RFDE32109876543210987654"},
			},
			wait: false,
		},
		{
			desc: "NOSuccess. Country already exists",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901234567890123456": countryAcc,
				},
			},
			countryAccount: countryAcc,
			destructionAccount: destructionAcc,
			wait: true,
		},
		{
			desc: "NOSuccess. Destruction already exists",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": destructionAcc,
				},
			},
			countryAccount: countryAcc,
			destructionAccount: destructionAcc,
			wait: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.db.AddGovernment(context.Background(), tC.countryAccount, tC.destructionAccount)
			if (err != nil) != tC.wait {
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.wait, err)
			}
		})
	}
}

func TestAddAccount(t *testing.T) {
	accou := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
	}
	testCases := []struct {
		desc string
		db   *repository.LocalDB
		acc  *account.Account
		wait bool
	}{
		{
			desc: "Success",
			db: &repository.LocalDB{
				Accounts: make(map[string]*account.Account),
			},
			acc: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
			},
			wait: false,
		},
		{
			desc: "NOSuccess. account already exists",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": accou,
				},
			},
			acc: accou,
			wait: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.db.AddAccount(context.Background(), tC.acc)
			if (err != nil) != tC.wait {
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.wait, err)
			}
		})
	}
}

func TestGetCountryIBAN(t *testing.T) {
	countryAcc := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
	}
	testCases := []struct {
		desc      string
		db        *repository.LocalDB
		waitValye tools.Iban
		waitErr   bool
	}{
		{
			desc: "Success",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": countryAcc,
				},
				CountryAccount: countryAcc,
			},
			waitValye: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
			waitErr:   false,
		},
		{
			desc: "Empty country account",
			db: &repository.LocalDB{
				Accounts: make(map[string]*account.Account),
			},
			waitValye: tools.Iban{},
			waitErr:   true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := tC.db.GetCountryIBAN(context.Background())
			if (err != nil) != tC.waitErr {
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.waitErr, err)
			}
		})
	}
}

func TestGetDestroyIBAN(t *testing.T) {
	destructionAcc := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
	}
	testCases := []struct {
		desc      string
		db        *repository.LocalDB
		waitValye tools.Iban
		waitErr   bool
	}{
		{
			desc: "Success",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": destructionAcc,
				},
				DestructionAccount: destructionAcc,
			},
			waitValye: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
			waitErr:   false,
		},
		{
			desc: "Empty country account",
			db: &repository.LocalDB{
				Accounts: make(map[string]*account.Account),
			},
			waitValye: tools.Iban{},
			waitErr:   true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := tC.db.GetDestroyIBAN(context.Background())
			if (err != nil) != tC.waitErr {
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.waitErr, err)
			}
		})
	}
}

func TestEmition(t *testing.T) {
	countryAcc := &account.Account{
		AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
		Balance: 100,
	}
	testCases := []struct {
		desc    string
		db      *repository.LocalDB
		amount  float64
		balanceAfter float64
		waitErr bool
	}{
		{
			desc: "Success",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": countryAcc,
				},
				CountryAccount: countryAcc,
			},
			amount: 100,
			balanceAfter: 200,
			waitErr:   false,
		},
		{
			desc: "Success2",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": countryAcc,
				},
				CountryAccount: countryAcc,
			},
			amount: 100,
			balanceAfter: 300,
			waitErr:   false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.db.Emition(context.Background(), tC.amount)
			if (err != nil) != tC.waitErr {
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.waitErr, err)
			}
			assert.Equal(t, tC.balanceAfter, tC.db.CountryAccount.Balance, "Incorrectly issued ")
		})
	}
}
