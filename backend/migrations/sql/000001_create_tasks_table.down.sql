-- トリガーの削除
DROP TRIGGER IF EXISTS tasks_updated_at_trigger ON tasks;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_updated_at_column();

-- タスクテーブルの削除
DROP TABLE IF EXISTS tasks;

-- ユーザーテーブルの削除
DROP TABLE IF EXISTS users;
