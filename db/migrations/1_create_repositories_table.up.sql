CREATE TABLE repositories (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    url TEXT NOT NULL UNIQUE,
    provider TEXT NOT NULL,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    slug TEXT NOT NULL,
    ref TEXT NOT NULL DEFAULT 'HEAD'
);

