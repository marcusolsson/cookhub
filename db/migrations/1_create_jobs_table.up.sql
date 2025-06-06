CREATE TYPE status_enum AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE jobs (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    slug TEXT NOT NULL,
    commit_sha TEXT NOT NULL UNIQUE,
    status STATUS_ENUM NOT NULL DEFAULT 'pending'
);
