package camunda

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Runtime struct {
	resty *resty.Client
}

func NewRuntime(client *resty.Client) *Runtime {
	return &Runtime{resty: client}
}

func (r *Runtime) StartProcess(ctx context.Context, flowID string, payload map[string]any) error {
	resp, err := r.resty.R().
		SetContext(ctx).
		SetBody(map[string]any{
			"variables": payload,
		}).
		Post(fmt.Sprintf("/process-definition/key/%s/start", flowID))
	if err != nil {
		return fmt.Errorf("start process: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("start process error: %s", resp.String())
	}
	return nil
}

func (r *Runtime) RetryProcess(ctx context.Context, workOrderID string) error {
	resp, err := r.resty.R().
		SetContext(ctx).
		Post(fmt.Sprintf("/external-task/%s/retry", workOrderID))
	if err != nil {
		return fmt.Errorf("retry process: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("retry process error: %s", resp.String())
	}
	return nil
}
