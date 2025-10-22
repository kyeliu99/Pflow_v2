package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/example/pflow/backend/internal/config"
	"github.com/example/pflow/backend/internal/flow"
	"github.com/example/pflow/backend/internal/workorder"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     config.QueueConfig
}

func NewPublisher(cfg config.QueueConfig) (*Publisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to queue: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &Publisher{conn: conn, channel: ch, cfg: cfg}, nil
}

func (p *Publisher) PublishFlowCreated(ctx context.Context, flow flow.Flow) error {
	return p.publish(ctx, "flow.created", flow)
}

func (p *Publisher) PublishFlowUpdated(ctx context.Context, flow flow.Flow) error {
	return p.publish(ctx, "flow.updated", flow)
}

func (p *Publisher) PublishWorkOrderCreated(ctx context.Context, wo workorder.WorkOrder) error {
	return p.publish(ctx, "workorder.created", wo)
}

func (p *Publisher) PublishWorkOrderCompleted(ctx context.Context, wo workorder.WorkOrder) error {
	return p.publish(ctx, "workorder.completed", wo)
}

func (p *Publisher) publish(ctx context.Context, event string, payload any) error {
	if p == nil || p.channel == nil {
		return nil
	}

	body, err := json.Marshal(struct {
		Event string `json:"event"`
		Data  any    `json:"data"`
	}{Event: event, Data: payload})
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return p.channel.Publish(p.cfg.Exchange, p.cfg.RoutingKey, false, false, amqp.Publishing{
		ContentType: p.cfg.ContentType,
		Body:        body,
	})
}

func (p *Publisher) Close() error {
	if p == nil {
		return nil
	}
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
