CREATE USER notes_user WITH PASSWORD 'notes_password';
CREATE DATABASE notes_db OWNER notes_user;

\connect notes_db;

CREATE TABLE IF NOT EXISTS notes (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
