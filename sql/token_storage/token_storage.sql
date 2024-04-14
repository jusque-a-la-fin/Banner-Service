-- название БД: token_storage. Хранит информацию о токенах
BEGIN;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  user_token VARCHAR(160) UNIQUE NOT NULL
);

CREATE TABLE admins (
  id SERIAL PRIMARY KEY,
  admin_token VARCHAR(160) UNIQUE NOT NULL
);

INSERT INTO users (user_token)
VALUES ('user_token');

INSERT INTO admins (admin_token)
VALUES ('admin_token');

COMMIT;
