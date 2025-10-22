package flow

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]Flow, error)
	Create(ctx context.Context, input CreateInput) (Flow, error)
	Get(ctx context.Context, id string) (Flow, error)
	Update(ctx context.Context, input UpdateInput) (Flow, error)
}

type Repository interface {
	List(ctx context.Context) ([]Flow, error)
	Get(ctx context.Context, id string) (Flow, error)
	Create(ctx context.Context, flow Flow) (Flow, error)
	Update(ctx context.Context, flow Flow) (Flow, error)
}

type CamundaDeployer interface {
	Deploy(ctx context.Context, flow Flow) error
}

type Publisher interface {
	PublishFlowCreated(ctx context.Context, flow Flow) error
	PublishFlowUpdated(ctx context.Context, flow Flow) error
}

type CreateInput struct {
	Name        string
	Description string
	Definition  map[string]any
	Metadata    map[string]string
}

type UpdateInput struct {
	ID          string
	Description string
	Definition  map[string]any
	Metadata    map[string]string
}

type service struct {
	repo      Repository
	camunda   CamundaDeployer
	publisher Publisher
}

type notFoundError struct{ id string }

func (e notFoundError) Error() string { return fmt.Sprintf("flow %s not found", e.id) }

func (notFoundError) NotFound() {}

func IsNotFound(err error) bool {
	var target notFoundError
	return errors.As(err, &target)
}

func NewService(repo Repository, camunda CamundaDeployer, publisher Publisher) Service {
	return &service{repo: repo, camunda: camunda, publisher: publisher}
}

func (s *service) List(ctx context.Context) ([]Flow, error) {
	return s.repo.List(ctx)
}

func (s *service) Create(ctx context.Context, input CreateInput) (Flow, error) {
	if input.Name == "" {
		return Flow{}, errors.New("name is required")
	}

	flow := Flow{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		Definition:  input.Definition,
		Metadata:    input.Metadata,
		Version:     1,
	}

	saved, err := s.repo.Create(ctx, flow)
	if err != nil {
		return Flow{}, err
	}

	if s.camunda != nil {
		if err := s.camunda.Deploy(ctx, saved); err != nil {
			return Flow{}, fmt.Errorf("deploy to camunda: %w", err)
		}
	}

	if s.publisher != nil {
		if err := s.publisher.PublishFlowCreated(ctx, saved); err != nil {
			return Flow{}, fmt.Errorf("publish flow created: %w", err)
		}
	}

	return saved, nil
}

func (s *service) Get(ctx context.Context, id string) (Flow, error) {
	flow, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sqlErrNotFound) {
			return Flow{}, notFoundError{id: id}
		}
		return Flow{}, err
	}
	return flow, nil
}

func (s *service) Update(ctx context.Context, input UpdateInput) (Flow, error) {
	existing, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		if errors.Is(err, sqlErrNotFound) {
			return Flow{}, notFoundError{id: input.ID}
		}
		return Flow{}, err
	}

	existing.Description = input.Description
	existing.Definition = input.Definition
	existing.Metadata = input.Metadata
	existing.Version++

	saved, err := s.repo.Update(ctx, existing)
	if err != nil {
		return Flow{}, err
	}

	if s.camunda != nil {
		if err := s.camunda.Deploy(ctx, saved); err != nil {
			return Flow{}, fmt.Errorf("deploy to camunda: %w", err)
		}
	}

	if s.publisher != nil {
		if err := s.publisher.PublishFlowUpdated(ctx, saved); err != nil {
			return Flow{}, fmt.Errorf("publish flow updated: %w", err)
		}
	}

	return saved, nil
}

var sqlErrNotFound = errors.New("flow not found")

func WrapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sqlErrNotFound) {
		return notFoundError{}
	}
	return err
}
