-- name: ListIngestionRuns :many
select * from ingestion_runs;

-- name: CreateIngestionRun :one
insert into ingestion_runs (repository_id, branch, commit_sha) values (
    $1, $2, $3
) returning id;

-- name: ListRepositories :many
select * from repositories;

-- name: ImportFiles :copyfrom
insert into files (
    ingestion_run_id, path, basename, stem, extension, content, size_bytes, hash
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: MarkRunAsCompleted :exec
update ingestion_runs set
    status = 'completed', files_processed_count = $1, completed_at = now()
where id = $2;

-- name: MarkRunAsFailed :exec
update ingestion_runs set
    status = 'failed',
    files_processed_count = $1,
    completed_at = now(),
    error_message = $2
where id = $3;

-- name: GetFile :one
WITH latest_run AS (
  SELECT ir.*,
         ROW_NUMBER() OVER (ORDER BY ir.started_at DESC) as rn
  FROM ingestion_runs ir
  JOIN repositories r ON ir.repository_id = r.id
  WHERE r.provider = $1 
    AND r.owner = $2 
    AND r.repo_name = $3
)
SELECT f.*
FROM files f
JOIN latest_run lr ON f.ingestion_run_id = lr.id
WHERE lr.rn = 1
  AND f.path = $4;

-- name: GetFiles :many
WITH latest_run AS (
  SELECT r.provider, r.owner, r.repo_name, ir.*,
         ROW_NUMBER() OVER (ORDER BY ir.started_at DESC) as rn
  FROM ingestion_runs ir
  JOIN repositories r ON ir.repository_id = r.id
  WHERE r.provider = $1 
    AND r.owner = $2 
    AND r.repo_name = $3
)
SELECT lr.provider, lr.owner, lr.repo_name, f.*
FROM files f
JOIN latest_run lr ON f.ingestion_run_id = lr.id
WHERE lr.rn = 1;
