package workorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]WorkOrder, error) {
	const query = `SELECT id, flow_id, title, assignee, status, payload, metadata, created_at, updated_at FROM workorders ORDER BY created_at DESC`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list workorders: %w", err)
	}
	defer rows.Close()

	var result []WorkOrder
	for rows.Next() {
		wo, err := scanWorkOrder(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, wo)
	}

	return result, rows.Err()
}

func (r *repository) Get(ctx context.Context, id string) (WorkOrder, error) {
	const query = `SELECT id, flow_id, title, assignee, status, payload, metadata, created_at, updated_at FROM workorders WHERE id = $1`

	row := r.db.QueryRowxContext(ctx, query, id)

	wo, err := scanWorkOrder(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkOrder{}, sqlErrNotFound
		}
		return WorkOrder{}, err
	}

	return wo, nil
}

func (r *repository) Create(ctx context.Context, wo WorkOrder) (WorkOrder, error) {
	const query = `INSERT INTO workorders (id, flow_id, title, assignee, status, payload, metadata, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	payload, err := json.Marshal(wo.Payload)
	if err != nil {
		return WorkOrder{}, fmt.Errorf("marshal payload: %w", err)
	}
	metadata, err := json.Marshal(wo.Metadata)
	if err != nil {
		return WorkOrder{}, fmt.Errorf("marshal metadata: %w", err)
	}

	now := time.Now().UTC()
	wo.CreatedAt = now
	wo.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query, wo.ID, wo.FlowID, wo.Title, wo.Assignee, wo.Status, payload, metadata, wo.CreatedAt, wo.UpdatedAt)
	if err != nil {
		return WorkOrder{}, fmt.Errorf("insert workorder: %w", err)
	}

	return wo, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status Status) error {
	const query = `UPDATE workorders SET status = $2, updated_at = $3 WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id, status, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return sqlErrNotFound
	}

	return nil
}

func scanWorkOrder(scanner interface {
	Scan(dest ...any) error
}) (WorkOrder, error) {
	var (
		wo          WorkOrder
		payloadRaw  []byte
		metadataRaw []byte
	)

	if err := scanner.Scan(&wo.ID, &wo.FlowID, &wo.Title, &wo.Assignee, &wo.Status, &payloadRaw, &metadataRaw, &wo.CreatedAt, &wo.UpdatedAt); err != nil {
		return WorkOrder{}, err
	}

	if len(payloadRaw) > 0 {
		if err := json.Unmarshal(payloadRaw, &wo.Payload); err != nil {
			return WorkOrder{}, fmt.Errorf("unmarshal payload: %w", err)
		}
	}

	if len(metadataRaw) > 0 {
		if err := json.Unmarshal(metadataRaw, &wo.Metadata); err != nil {
			return WorkOrder{}, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	return wo, nil
}
