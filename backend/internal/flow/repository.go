package flow

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/kyeliu99/Pflow_v2/backend/internal/persistence"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]Flow, error) {
	const query = `SELECT id, name, description, definition, metadata, version, created_at, updated_at FROM flows ORDER BY updated_at DESC`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list flows: %w", err)
	}
	defer rows.Close()

	var flows []Flow
	for rows.Next() {
		var (
			f           Flow
			definition  []byte
			metadataRaw []byte
		)
		if err := rows.Scan(&f.ID, &f.Name, &f.Description, &definition, &metadataRaw, &f.Version, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan flow: %w", err)
		}
		if len(definition) > 0 {
			if err := json.Unmarshal(definition, &f.Definition); err != nil {
				return nil, fmt.Errorf("unmarshal definition: %w", err)
			}
		}
		if len(metadataRaw) > 0 {
			if err := json.Unmarshal(metadataRaw, &f.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal metadata: %w", err)
			}
		}
		flows = append(flows, f)
	}

	return flows, rows.Err()
}

func (r *repository) Get(ctx context.Context, id string) (Flow, error) {
	const query = `SELECT id, name, description, definition, metadata, version, created_at, updated_at FROM flows WHERE id = $1`

	var (
		f           Flow
		definition  []byte
		metadataRaw []byte
	)

	err := r.db.QueryRowxContext(ctx, query, id).Scan(&f.ID, &f.Name, &f.Description, &definition, &metadataRaw, &f.Version, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, persistence.ErrNotFound) {
			return Flow{}, sqlErrNotFound
		}
		return Flow{}, fmt.Errorf("get flow: %w", err)
	}

	if len(definition) > 0 {
		if err := json.Unmarshal(definition, &f.Definition); err != nil {
			return Flow{}, fmt.Errorf("unmarshal definition: %w", err)
		}
	}

	if len(metadataRaw) > 0 {
		if err := json.Unmarshal(metadataRaw, &f.Metadata); err != nil {
			return Flow{}, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	return f, nil
}

func (r *repository) Create(ctx context.Context, flow Flow) (Flow, error) {
	const query = `INSERT INTO flows (id, name, description, definition, metadata, version, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	definition, err := json.Marshal(flow.Definition)
	if err != nil {
		return Flow{}, fmt.Errorf("marshal definition: %w", err)
	}

	metadata, err := json.Marshal(flow.Metadata)
	if err != nil {
		return Flow{}, fmt.Errorf("marshal metadata: %w", err)
	}

	now := time.Now().UTC()
	flow.CreatedAt = now
	flow.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query, flow.ID, flow.Name, flow.Description, definition, metadata, flow.Version, flow.CreatedAt, flow.UpdatedAt)
	if err != nil {
		return Flow{}, fmt.Errorf("insert flow: %w", err)
	}

	return flow, nil
}

func (r *repository) Update(ctx context.Context, flow Flow) (Flow, error) {
	const query = `UPDATE flows SET description = $2, definition = $3, metadata = $4, version = $5, updated_at = $6 WHERE id = $1`

	definition, err := json.Marshal(flow.Definition)
	if err != nil {
		return Flow{}, fmt.Errorf("marshal definition: %w", err)
	}

	metadata, err := json.Marshal(flow.Metadata)
	if err != nil {
		return Flow{}, fmt.Errorf("marshal metadata: %w", err)
	}

	flow.UpdatedAt = time.Now().UTC()

	res, err := r.db.ExecContext(ctx, query, flow.ID, flow.Description, definition, metadata, flow.Version, flow.UpdatedAt)
	if err != nil {
		return Flow{}, fmt.Errorf("update flow: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return Flow{}, fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return Flow{}, sqlErrNotFound
	}

	return flow, nil
}
