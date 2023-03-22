package account

import (
	"context"
	"fmt"

	"encore.app/account/workflow"
	"encore.app/domain"
	"encore.dev"
	"encore.dev/rlog"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Use an environment-specific task queue so we can use the same
// Temporal Cluster for all cloud environments.
var (
	envName          = encore.Meta().Environment.Name
	accountTaskQueue = envName + "-greeting"
)

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}

//encore:api auth method=POST path=/account/create
func (ac *Service) CreateAccount(ctx context.Context, account domain.Account) (*domain.Response, error) {
	options := client.StartWorkflowOptions{
		ID:        "create-account-workflow",
		TaskQueue: accountTaskQueue + "-taskqueue",
	}
	we, err := ac.client.ExecuteWorkflow(ctx, options, workflow.CreateAccount, ac.db, account)
	if err != nil {
		return nil, err
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	// Get the results
	var response domain.Response
	err = we.Get(ctx, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//encore:api auth method=GET path=/account/:id
func (ac *Service) GetAccount(ctx context.Context, id string) (*domain.Response, error) {
	options := client.StartWorkflowOptions{
		ID:        "get-account-workflow",
		TaskQueue: accountTaskQueue + "-taskqueue",
	}
	we, err := ac.client.ExecuteWorkflow(ctx, options, workflow.GetAccount, ac.db, id)
	if err != nil {
		return nil, err
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	// Get the results
	var response domain.Response
	err = we.Get(ctx, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//encore:api auth method=POST path=/account/transfer
func (ac *Service) TransferAccount(ctx context.Context, account *domain.TransferAmount) (*domain.Response, error) {
	options := client.StartWorkflowOptions{
		ID:        "amount-transfer-workflow",
		TaskQueue: accountTaskQueue + "-taskqueue",
	}
	we, err := ac.client.ExecuteWorkflow(ctx, options, workflow.TransferAccount, ac.db, account)
	if err != nil {
		return nil, err
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	// Get the results
	var response domain.Response
	err = we.Get(ctx, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//encore:service
type Service struct {
	db     tb.Client
	client client.Client
	worker worker.Worker
}

func initService() (*Service, error) {
	db, err := tb.NewClient(0, []string{"3000"}, 1)
	if err != nil {
		fmt.Println("error on creation: " + err.Error())
		return nil, err
	}

	c, err := client.Dial(client.Options{HostPort: "127.0.0.1:7233"})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, "tms-app", worker.Options{})

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}

	w.RegisterWorkflow(workflow.CreateAccount)
	w.RegisterActivity(workflow.CreateAccountActivity)
	w.RegisterWorkflow(workflow.GetAccount)
	w.RegisterActivity(workflow.GetAccountActivity)
	w.RegisterWorkflow(workflow.TransferAccount)
	w.RegisterActivity(workflow.TransferAccountActivity)

	return &Service{db: db, client: c, worker: w}, nil
}
