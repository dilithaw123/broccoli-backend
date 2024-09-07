CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE NOT NULL
);

CREATE TABLE groups (
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  allowed_emails TEXT[]
);

CREATE TABLE sessions (
  id BIGSERIAL PRIMARY KEY,
  group_id BIGINT REFERENCES groups(id),
  create_date DATE,
  UNIQUE (group_id, create_date)
);

CREATE TABLE user_submissions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id),
  session_id BIGINT REFERENCES sessions(id),
  yesterday TEXT[],
  today TEXT[],
  blockers TEXT[]
);
