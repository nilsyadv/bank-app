package workflow

import (
	"time"

	"encore.app/domain"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	"go.temporal.io/sdk/workflow"
)

func CreateAccount(ctx workflow.Context, db tb.Client, acc domain.Account) (*domain.Response, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var result domain.Response
	err := workflow.ExecuteActivity(ctx, CreateAccountActivity, db, acc).Get(ctx, &result)
	return &result, err
}

func TransferAccount(ctx workflow.Context, db tb.Client, trfaccount domain.TransferAmount) (*domain.Response, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var result domain.Response
	err := workflow.ExecuteActivity(ctx, TransferAccountActivity, db, trfaccount).Get(ctx, &result)
	return &result, err
}

func GetAccount(ctx workflow.Context, db tb.Client, accountid string) (*domain.Response, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var result domain.Response
	err := workflow.ExecuteActivity(ctx, GetAccountActivity, db, accountid).Get(ctx, &result)
	return &result, err
}
