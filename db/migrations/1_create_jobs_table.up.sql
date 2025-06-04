CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    slug TEXT NOT NULL,
    commit_sha TEXT NOT NULL UNIQUE
)
