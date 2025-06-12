CREATE TYPE status AS ENUM ('processing', 'completed', 'failed');

CREATE TABLE ingestion_runs (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    repository_id TEXT NOT NULL,
    branch TEXT NOT NULL,
    commit_sha TEXT NOT NULL,
    status STATUS NOT NULL DEFAULT 'processing',
    started_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    files_processed_count INTEGER,
    completed_at TIMESTAMPTZ,
    error_message TEXT,

    FOREIGN KEY (repository_id) REFERENCES repositories (id)
);
