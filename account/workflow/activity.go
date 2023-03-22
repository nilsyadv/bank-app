package workflow

import (
	"errors"
	"log"

	"encore.app/domain"
	"encore.app/utility"
	"encore.dev/beta/errs"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	"github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type Account struct {
	ID     string `json:"id,omitempty"`
	Ledger uint32 `json:"ledger,omitempty"`
	Code   uint16 `json:"code,omitempty"`
}

func CreateAccountActivity(db tb.Client, acc domain.Account) (*domain.Response, error) {
	res, err := db.CreateAccounts([]types.Account{
		{
			ID:     utility.Uint128(acc.ID),
			Ledger: acc.Ledger,
			Code:   acc.Code,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, err := range res {
		return nil, errors.New(err.Result.String())
	}

	return &domain.Response{Message: "account created successfully"}, nil
}

func GetAccountActivity(db tb.Client, accountid string) ([]domain.Account, error) {
	dbaccounts, err := db.LookupAccounts([]types.Uint128{
		utility.Uint128(accountid),
	})
	if err != nil {
		log.Printf("Error creating accounts: %s", err)
		return nil, err
	}

	var accounts []domain.Account
	for _, dbaccount := range dbaccounts {
		accounts = append(accounts, domain.Account{
			ID:     dbaccount.ID.String(),
			Ledger: dbaccount.Ledger,
			Code:   dbaccount.Code,
		})
	}

	if len(accounts) == 0 {
		return nil, &errs.Error{Message: "account not found", Code: 5}
	}

	return accounts, nil
}

func TransferAccountActivity(db tb.Client, trfaccount *domain.TransferAmount) (*domain.Response, error) {
	debitaccount, err := GetAccountActivity(db, trfaccount.SenderAccountID)
	if err != nil {
		return nil, err
	}
	creditaccount, err := GetAccountActivity(db, trfaccount.RecipientAccountID)
	if err != nil {
		return nil, err
	}
	if debitaccount[0].Ledger != creditaccount[0].Ledger {
		return nil, errors.New(types.AccountExistsWithDifferentLedger.String())
	}

	batch := types.Transfer{
		ID:              utility.NewTransactionID(),
		DebitAccountID:  utility.Uint128(debitaccount[0].ID),
		CreditAccountID: utility.Uint128(creditaccount[0].ID),
		Ledger:          debitaccount[0].Ledger,
		Code:            debitaccount[0].Code,
		Amount:          uint64(trfaccount.Amount),
	}

	res, err := db.CreateTransfers([]types.Transfer{batch})
	if err != nil {
		return nil, err
	}
	for _, err := range res {
		return nil, &errs.Error{
			Code:    10,
			Message: "failed to transfer amount to ledger: " + err.Result.String(),
		}
	}
	return &domain.Response{Message: "amount transfer successfully"}, nil
}
