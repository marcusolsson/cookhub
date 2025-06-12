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
with latest_run as (
    select
        ir.*,
        row_number() over (
            order by ir.started_at desc
        ) as rn
    from ingestion_runs ir
    join repositories r on ir.repository_id = r.id
    where
        r.provider = $1
        and r.owner = $2
        and r.repo_name = $3
)

select f.*
from files f
join latest_run lr on f.ingestion_run_id = lr.id
where
    lr.rn = 1
    and f.path = $4;

-- name: GetFiles :many
with latest_run as (
    select
        r.provider,
        r.owner,
        r.repo_name,
        ir.*,
        row_number() over (
            order by ir.started_at desc
        ) as rn
    from ingestion_runs ir
    join repositories r on ir.repository_id = r.id
    where
        r.provider = $1
        and r.owner = $2
        and r.repo_name = $3
)

select
    lr.provider,
    lr.owner,
    lr.repo_name,
    f.*
from files f
join latest_run lr on f.ingestion_run_id = lr.id
where lr.rn = 1;

-- name: GetCommitShaByRepo :one
with latest_run as (
    select
        ir.commit_sha,
        row_number() over (
            order by ir.started_at desc
        ) as rn
    from ingestion_runs ir
    join repositories r on ir.repository_id = r.id
    where r.id = $1
)

select commit_sha from latest_run
where rn = 1;
