CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    tags TEXT[],
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL
)
