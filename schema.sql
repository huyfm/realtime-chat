CREATE TABLE users (
    id        SERIAL PRIMARY KEY,
    name      TEXT NOT NULL,
    email     TEXT UNIQUE,
    github_id BIGINT UNIQUE NOT NULL
);
