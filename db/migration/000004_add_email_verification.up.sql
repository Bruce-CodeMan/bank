-- modify "user" table
ALTER TABLE "user"
ADD COLUMN "is_email_verified" boolean NOT NULL DEFAULT false;

-- create "verify_emails" table
CREATE TABLE "verify_emails" (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "user_public_id" uuid NOT NULL,
    "email" varchar NOT NULL, 
    "secret_code" varchar NOT NULL,
    "is_used" boolean NOT NULL DEFAULT false,
    "expired_at" timestamptz NOT NULL DEFAULT (now() + INTERVAL '15 minutes'),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- add index
CREATE INDEX idx_verify_emails_user_public_id ON verify_emails(user_public_id);
CREATE INDEX idx_verify_emails_secret_code ON verify_emails(secret_code);
CREATE INDEX idx_verify_emails_expires_at ON verify_emails(expired_at);