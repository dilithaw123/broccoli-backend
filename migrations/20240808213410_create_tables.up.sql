CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE groups (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  allowed_emails TEXT[],
  timezone TEXT NOT NULL
);

CREATE TABLE sessions (
  id BIGSERIAL PRIMARY KEY,
  group_id BIGINT REFERENCES groups(id),
  create_date TIMESTAMP WITH TIME ZONE NOT NULL,
);

CREATE TABLE user_submissions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  session_id BIGINT REFERENCES sessions(id),
  yesterday TEXT[],
  today TEXT[],
  blockers TEXT[]
);
