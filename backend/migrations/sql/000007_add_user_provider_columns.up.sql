-- ユーザーテーブルにプロバイダー関連カラムを追加
-- 007_add_user_provider_columns.up.sql

-- Email カラム追加（UNIQUE制約付き）
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- Provider カラム追加（デフォルト値: Cognito_UserPool）
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50) DEFAULT 'Cognito_UserPool';

-- Provider ID カラム追加（Google sub等を格納、NULL許可）
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider_id VARCHAR(255);

-- Updated At カラム追加
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Email ユニーク制約追加（NULL値は除外）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'users_email_key'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
    END IF;
END $$;

-- Provider + Provider ID の組み合わせユニーク制約（同一プロバイダー内でのID重複防止）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'users_provider_unique'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT users_provider_unique UNIQUE (provider, provider_id);
    END IF;
END $$;

-- Updated At トリガー追加（既存のupdate_updated_at_column関数を使用）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger 
        WHERE tgname = 'users_updated_at_trigger'
    ) THEN
        CREATE TRIGGER users_updated_at_trigger
        BEFORE UPDATE ON users
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- 既存データの更新（テストユーザーのemail設定）
UPDATE users 
SET email = 'dev@example.com', 
    updated_at = NOW()
WHERE id = '00000000-0000-0000-0000-000000000001' 
  AND email IS NULL;

-- インデックス追加（パフォーマンス向上）
-- CONCURRENTLYを削除（トランザクション内では使用不可）
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider);
CREATE INDEX IF NOT EXISTS idx_users_provider_id ON users(provider_id) WHERE provider_id IS NOT NULL;
