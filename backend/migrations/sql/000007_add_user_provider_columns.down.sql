-- ユーザーテーブルからプロバイダー関連カラムを削除
-- 007_add_user_provider_columns.down.sql

-- トリガーの削除
DROP TRIGGER IF EXISTS users_updated_at_trigger ON users;

-- インデックスの削除
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_provider;
DROP INDEX IF EXISTS idx_users_provider_id;

-- 制約の削除
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_provider_unique;

-- カラムの削除
ALTER TABLE users DROP COLUMN IF EXISTS email;
ALTER TABLE users DROP COLUMN IF EXISTS provider;
ALTER TABLE users DROP COLUMN IF EXISTS provider_id;
ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
