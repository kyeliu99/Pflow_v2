CREATE TABLE IF NOT EXISTS flows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    definition JSONB NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS workorders (
    id TEXT PRIMARY KEY,
    flow_id TEXT NOT NULL REFERENCES flows(id),
    title TEXT NOT NULL,
    assignee TEXT,
    status TEXT NOT NULL,
    payload JSONB DEFAULT '{}'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_workorders_flow_id ON workorders(flow_id);
CREATE INDEX IF NOT EXISTS idx_workorders_status ON workorders(status);
