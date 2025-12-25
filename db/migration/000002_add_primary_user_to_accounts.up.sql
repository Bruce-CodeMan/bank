CREATE TABLE "user" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "public_id" uuid UNIQUE NOT NULL,
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- account新增primary_user_id(先允许为 NULL)
ALTER TABLE account 
ADD COLUMN primary_user_id bigint;

-- 回填 primary_user_id
UPDATE account a 
SET primary_user_id = u.id
FROM "user" u
WHERE a.owner = u.username;

-- 加 NOT NULL
ALTER TABLE account 
ALTER COLUMN primary_user_id SET NOT NULL;

-- 加唯一约束
CREATE UNIQUE INDEX uniq_account_primary_user_currency
ON account(primary_user_id, currency);

-- 创建account_users表
CREATE TABLE "account_users" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- account_users添加唯一约束
CREATE UNIQUE INDEX uniq_account_users
ON account_users(account_id, user_id);

-- 把owner写入account_users
INSERT INTO account_users(account_id, user_id, role)
SELECT a.id, a.primary_user_id, 'owner'
FROM account a;

-- 删除旧字段
ALTER TABLE account DROP COLUMN owner;