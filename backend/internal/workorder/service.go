package workorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]WorkOrder, error)
	Create(ctx context.Context, input CreateInput) (WorkOrder, error)
	Get(ctx context.Context, id string) (WorkOrder, error)
	Retry(ctx context.Context, id string) error
}

type Repository interface {
	List(ctx context.Context) ([]WorkOrder, error)
	Get(ctx context.Context, id string) (WorkOrder, error)
	Create(ctx context.Context, wo WorkOrder) (WorkOrder, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
}

type FlowReader interface {
	Get(ctx context.Context, id string) (FlowSummary, error)
}

type FlowSummary struct {
	ID   string
	Name string
}

type CamundaRuntime interface {
	StartProcess(ctx context.Context, flowID string, payload map[string]any) error
	RetryProcess(ctx context.Context, workOrderID string) error
}

type Publisher interface {
	PublishWorkOrderCreated(ctx context.Context, wo WorkOrder) error
	PublishWorkOrderCompleted(ctx context.Context, wo WorkOrder) error
}

type CreateInput struct {
	FlowID   string
	Title    string
	Assignee string
	Payload  map[string]any
	Metadata map[string]string
}

type service struct {
	repo      Repository
	flows     FlowReader
	runtime   CamundaRuntime
	publisher Publisher
}

type notFoundError struct{ id string }

func (e notFoundError) Error() string { return fmt.Sprintf("workorder %s not found", e.id) }

func (notFoundError) NotFound() {}

func IsNotFound(err error) bool {
	var target notFoundError
	return errors.As(err, &target)
}

func NewService(repo Repository, flows FlowReader, runtime CamundaRuntime, publisher Publisher) Service {
	return &service{repo: repo, flows: flows, runtime: runtime, publisher: publisher}
}

func (s *service) List(ctx context.Context) ([]WorkOrder, error) {
	return s.repo.List(ctx)
}

func (s *service) Create(ctx context.Context, input CreateInput) (WorkOrder, error) {
	if input.FlowID == "" {
		return WorkOrder{}, errors.New("flow id is required")
	}
	if input.Title == "" {
		return WorkOrder{}, errors.New("title is required")
	}

	flow, err := s.flows.Get(ctx, input.FlowID)
	if err != nil {
		return WorkOrder{}, fmt.Errorf("load flow: %w", err)
	}

	wo := WorkOrder{
		ID:       uuid.NewString(),
		FlowID:   flow.ID,
		Title:    input.Title,
		Assignee: input.Assignee,
		Status:   StatusPending,
		Payload:  input.Payload,
		Metadata: input.Metadata,
	}

	saved, err := s.repo.Create(ctx, wo)
	if err != nil {
		return WorkOrder{}, err
	}

	if s.runtime != nil {
		if err := s.runtime.StartProcess(ctx, saved.FlowID, saved.Payload); err != nil {
			return WorkOrder{}, fmt.Errorf("start process: %w", err)
		}
	}

	if s.publisher != nil {
		if err := s.publisher.PublishWorkOrderCreated(ctx, saved); err != nil {
			return WorkOrder{}, fmt.Errorf("publish workorder: %w", err)
		}
	}

	return saved, nil
}

func (s *service) Get(ctx context.Context, id string) (WorkOrder, error) {
	wo, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sqlErrNotFound) {
			return WorkOrder{}, notFoundError{id: id}
		}
		return WorkOrder{}, err
	}
	return wo, nil
}

func (s *service) Retry(ctx context.Context, id string) error {
	wo, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sqlErrNotFound) {
			return notFoundError{id: id}
		}
		return err
	}

	if s.runtime != nil {
		if err := s.runtime.RetryProcess(ctx, wo.ID); err != nil {
			return fmt.Errorf("retry process: %w", err)
		}
	}

	return s.repo.UpdateStatus(ctx, id, StatusRunning)
}

var sqlErrNotFound = errors.New("workorder not found")
