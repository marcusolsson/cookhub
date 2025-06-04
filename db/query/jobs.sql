-- name: CreateJob :one
INSERT INTO jobs (slug, commit_sha) VALUES ($1, $2) RETURNING id;

-- name: CreateFile :exec
INSERT INTO files (job_id, name, content) VALUES ($1, $2, $3);

-- name: GetLatestJobBySlug :one
SELECT id FROM jobs
WHERE slug = $1
ORDER BY created_at DESC LIMIT 1;

-- name: GetFilesByJob :many
SELECT
    created_at,
    name,
    content
FROM files
WHERE job_id = $1;
