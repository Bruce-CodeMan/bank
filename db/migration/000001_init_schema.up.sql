CREATE TABLE "account" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "public_id" uuid UNIQUE NOT NULL,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entry" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "transfer_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfer" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())

  CONSTRAINT transfer_amount_positive CHECK (amount > 0)
  CONSTRAINT transfer_from_to_diff CHECK (from_account_id <> to_account_id)
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entry" ("account_id");

CREATE INDEX ON "transfer" ("from_account_id");

CREATE INDEX ON "transfer" ("to_account_id");

CREATE INDEX ON "transfer" ("from_account_id", "to_account_id");


