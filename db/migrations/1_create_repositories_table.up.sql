CREATE TABLE repositories (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    url TEXT NOT NULL,
    provider TEXT NOT NULL,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    branch TEXT NOT NULL
);

INSERT INTO repositories
(url, provider, owner, repo_name, branch)
VALUES
('github.com/cooklang/recipes', 'github.com', 'cooklang', 'recipes', 'HEAD'),
('github.com/marcusolsson/recipes', 'github.com', 'marcusolsson', 'recipes', 'main');
