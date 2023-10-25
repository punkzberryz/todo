CREATE TABLE "password_reset_sessions" (
    "email" varchar PRIMARY KEY,
    "otp" varchar NOT NULL,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "password_reset_sessions" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");
