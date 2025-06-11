-- name: CreateJob :one
INSERT INTO jobs (slug, commit_sha) VALUES ($1, $2) RETURNING id;

-- name: CreateFile :copyfrom
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

-- name: GetFileByName :one
SELECT
    created_at,
    name,
    content
FROM files
WHERE job_id = $1 AND name = $2;

-- name: GetJobs :many
SELECT * FROM jobs;

-- name: GetRecipes :many
SELECT
    f.created_at,
    f.name,
    f.content,
    j.slug,
    j.commit_sha
FROM files AS f
INNER JOIN jobs AS j ON f.job_id = j.id;

-- name: SetJobStatus :exec
UPDATE jobs SET status = $1
WHERE id = $2;

-- name: GetAllFiles :many
SELECT * FROM files;

-- name: AddRepositoryMetadata :exec
INSERT INTO repo_metadata (job_id, provider, owner, name, response) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetRepoMetadataByJob :one
SELECT * FROM repo_metadata
WHERE job_id = $1;

-- name: GetRecipePageData :one
SELECT
    f.content,
    rm.response
FROM files f
INNER JOIN
    repo_metadata rm
    ON f.job_id = rm.job_id AND f.job_id = $1 AND f.name = $2;

-- name: GetFilesByRepo :many
SELECT
    f.id,
    f.name,
    f.content,
    f.created_at,
    j.slug,
    j.commit_sha,
    j.status,
    rm.provider,
    rm.owner,
    rm.name AS repo_name
FROM files f
JOIN jobs j ON f.job_id = j.id
JOIN repo_metadata rm ON j.id = rm.job_id
WHERE
    j.id = (
        SELECT job_id
        FROM repo_metadata m
        WHERE m.provider = $1 AND m.owner = $2 AND m.name = $3
        ORDER BY m.created_at DESC
        LIMIT 1
    )
ORDER BY f.name;

-- name: GetAllRepos :many
SELECT DISTINCT ON (provider, owner, name)
    provider,
    owner,
    name
FROM repo_metadata
ORDER BY provider, owner, name, created_at DESC;

---
