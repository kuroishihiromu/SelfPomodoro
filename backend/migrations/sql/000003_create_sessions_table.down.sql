-- セッションテーブルの削除
DROP TRIGGER IF EXISTS sessions_updated_at_trigger ON sessions;
DROP TABLE IF EXISTS sessions;
