CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    job_id TEXT NOT NULL,
    name TEXT NOT NULL,
    content TEXT NOT NULL,

    FOREIGN KEY (job_id) REFERENCES jobs (id)
)
