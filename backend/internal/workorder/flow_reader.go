package workorder

import (
	"context"

	"github.com/kyeliu99/Pflow_v2/backend/internal/flow"
)

type FlowServiceAdapter struct {
	Service flow.Service
}

func (a FlowServiceAdapter) Get(ctx context.Context, id string) (FlowSummary, error) {
	f, err := a.Service.Get(ctx, id)
	if err != nil {
		return FlowSummary{}, err
	}

	return FlowSummary{ID: f.ID, Name: f.Name}, nil
}
