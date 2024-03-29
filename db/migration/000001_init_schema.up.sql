CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "user_role" varchar NOT NULL DEFAULT 'admin',
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
  "created_at" timestamptz NOT NULL DEFAULT (now())
);