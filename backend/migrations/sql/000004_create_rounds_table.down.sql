-- ラウンドテーブルの削除
DROP TRIGGER IF EXISTS rounds_updated_at_trigger ON rounds;
DROP TABLE IF EXISTS rounds;
