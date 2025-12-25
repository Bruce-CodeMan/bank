-- 1. account重新加回owner字段
ALTER TABLE account
ADD COLUMN owner varchar;

-- 2. 从primary_user_id 反向回填owner
UPDATE account a 
SET owner = u.username 
FROM "user" u 
WHERE a.primary_user_id = u.id;

-- 3. 删除account_users表
DROP TABLE IF EXISTS "account_users";

-- 4. 删除唯一索引
DROP INDEX IF EXISTS uniq_account_primary_user_currency;

-- 5. 删除primary_user_id字段
ALTER TABLE account 
DROP COLUMN IF EXISTS primary_user_id;

-- 6. 删除user表
DROP TABLE IF EXISTS "user";