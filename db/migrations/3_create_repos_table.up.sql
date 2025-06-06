CREATE TABLE repo_metadata (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    job_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    owner TEXT NOT NULL,
    name TEXT NOT NULL,
    response JSONB NOT NULL,

    FOREIGN KEY (job_id) REFERENCES jobs (id)
);
