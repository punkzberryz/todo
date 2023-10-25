CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email"   varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "body" varchar NOT NULL,
  "is_done" boolean NOT NULL DEFAULT false,
  "owner_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "tasks" ("owner_id");
CREATE INDEX ON "users" ("email");

ALTER TABLE "tasks" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");