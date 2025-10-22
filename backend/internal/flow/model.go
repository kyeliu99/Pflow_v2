package flow

import "time"

type Flow struct {
	ID          string            `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Definition  map[string]any    `json:"definition" db:"definition"`
	Metadata    map[string]string `json:"metadata" db:"metadata"`
	Version     int               `json:"version" db:"version"`
	CreatedAt   time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time         `json:"updatedAt" db:"updated_at"`
}

func (f Flow) Summary() map[string]any {
	return map[string]any{
		"id":          f.ID,
		"name":        f.Name,
		"description": f.Description,
		"version":     f.Version,
		"updatedAt":   f.UpdatedAt,
	}
}
