# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

This project uses [Task](https://taskfile.dev/) for build automation:

- `task serve` - Build and run local web server (uses .env.development)
- `task test` - Run all tests (`go test ./...`)
- `task generate` - Generate all code (runs sqlc and templ)
- `task templ` - Generate Templ templates only
- `task sqlc` - Generate Go types from SQL queries
- `task lint` - Run sqlfluff linter on SQL files
- `task format` - Format Go code and SQL files
- `task migrate` - Run database migrations (uses .env.production)

## Architecture

CookHub is a recipe aggregation platform that indexes CookLang recipes from GitHub repositories.

**Tech Stack:** Go 1.23.2, Chi (router), Templ (templates), SQLc (database), PostgreSQL, CookLang

**Request Flow:**
1. HTTP handlers in root (`pages.go`, `api.go`) receive requests
2. Database queries via SQLc-generated code in `db/sqlc/`
3. CookLang recipes parsed with `cooklang-go` library
4. HTML rendered with Templ components in `views/`

**Key Directories:**
- `db/query/` - SQL query definitions (edit these, then run `task sqlc`)
- `db/migrations/` - Database schema migrations
- `db/sqlc/` - Generated Go code (do not edit directly)
- `views/` - Templ templates (`.templ` files)
- `github/` - GitHub API client for repository sync
- `utils/` - Utility functions for recipe parsing

**Data Model:**
- `repositories` - Tracked GitHub repos containing recipes
- `ingestion_runs` - Import/sync session tracking
- `files` - Recipe file content and metadata

**Import Flow:** Repos are added via `POST /api/repos`, then `POST /api/import` downloads zipballs, extracts `.cook` files, and stores them with SHA256 hashes.
