package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (db *Store) CreateJobsTable(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		slug TEXT NOT NULL,
		commit_sha TEXT NOT NULL UNIQUE
	)`)
	return err
}

func (db *Store) CreateFilesTable(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		job_id TEXT NOT NULL,
    	name TEXT NOT NULL,
		content TEXT NOT NULL,

		FOREIGN KEY(job_id) REFERENCES jobs(id)
	)`)

	return err
}

func (s *Store) CreateJob(ctx context.Context, slug, commitSHA string) (string, error) {
	query := `
	INSERT INTO jobs (slug, commit_sha)
	VALUES ($1, $2)
	RETURNING id`

	var jobID string
	if err := s.pool.QueryRow(ctx, query,
		slug, commitSHA,
	).Scan(&jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *Store) CreateFile(ctx context.Context, jobID, name string, content []byte) error {
	query := `
	INSERT INTO files (job_id, name, content)
	VALUES ($1, $2, $3)
	`

	_, err := s.pool.Exec(ctx, query,
		jobID, name, string(content),
	)

	return err
}

func (s *Store) GetLatestJobBySlug(ctx context.Context, slug string) (string, error) {
	query := `
	SELECT id
	FROM jobs
	WHERE slug = $1
	ORDER BY created_at DESC
	LIMIT 1
	`
	var jobID string
	if err := s.pool.QueryRow(ctx, query, slug).Scan(&jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

type CooklangFile struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
}

func (s *Store) GetFilesByJob(ctx context.Context, jobID string) ([]CooklangFile, error) {
	query := `
	SELECT created_at, name, content
	FROM files
	WHERE job_id = $1
	`
	rows, err := s.pool.Query(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CooklangFile
	for rows.Next() {
		var entry CooklangFile
		if err := rows.Scan(&entry.CreatedAt, &entry.Name, &entry.Content); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
