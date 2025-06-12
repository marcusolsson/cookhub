CREATE TABLE files (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    ingestion_run_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    path TEXT NOT NULL,
    basename TEXT NOT NULL,
    stem TEXT NOT NULL,
    extension TEXT NOT NULL,
    content TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    hash BYTEA NOT NULL,

    FOREIGN KEY (ingestion_run_id) REFERENCES ingestion_runs (id)
);
