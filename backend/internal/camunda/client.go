package camunda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/kyeliu99/Pflow_v2/backend/internal/config"
	"github.com/kyeliu99/Pflow_v2/backend/internal/flow"
)

type Client struct {
	resty *resty.Client
}

func NewClient(cfg config.CamundaConfig) *Client {
	return &Client{resty: newResty(cfg)}
}

func newResty(cfg config.CamundaConfig) *resty.Client {
	return resty.New().
		SetBaseURL(cfg.BaseURL).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetHeader("Content-Type", "application/json")
}

func (c *Client) HTTP() *resty.Client {
	return c.resty
}

func (c *Client) Deploy(ctx context.Context, f flow.Flow) error {
	payload := map[string]any{
		"deployment-name":     fmt.Sprintf("pflow-%s", f.ID),
		"deploy-changed-only": true,
		"resources": map[string]string{
			fmt.Sprintf("%s.bpmn", f.ID): mustJSON(f.Definition),
		},
	}

	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/deployment/create")
	if err != nil {
		return fmt.Errorf("call camunda: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("camunda error: %s", resp.String())
	}

	return nil
}

func mustJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
