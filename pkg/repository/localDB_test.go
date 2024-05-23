package repository_test

import (
	"context"
	"testbanc/pkg/account"
	"testbanc/pkg/repository"
	"testbanc/pkg/tools"
	"testing"
)



func TestAddGovernment(t *testing.T) {
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
					"BY12UJIK78901234567890123456": &account.Account{
						AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
					},
				},
			},
			countryAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
			},
			destructionAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
			},
			wait: true,
		},
		{
			desc: "NOSuccess. Destruction already exists",
			db: &repository.LocalDB{
				Accounts: map[string]*account.Account{
					"BY12UJIK78901384957285950213": &account.Account{
						AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
					},
				},
			},
			countryAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901234567890123456"},
			},
			destructionAccount: &account.Account{
				AccountNumber: tools.Iban{Iban: "BY12UJIK78901384957285950213"},
			},
			wait: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.db.AddGovernment(context.Background(), tC.countryAccount, tC.destructionAccount)
			if (err != nil) != tC.wait{
				t.Errorf("expected an error indicator: %v, but the program broke and got an error: %v", tC.wait, err)
			}
		})
	}
}
