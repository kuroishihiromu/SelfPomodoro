-- テスト用ユーザーの追加
INSERT INTO users (id, name, created_at)
VALUES 
  ('00000000-0000-0000-0000-000000000001', 'テストユーザー', NOW())
ON CONFLICT (id) DO NOTHING;
