DROP TABLE IF EXISTS verify_emails;

ALTER TABLE "user"
DROP COLUMN IF EXISTS is_email_verified;