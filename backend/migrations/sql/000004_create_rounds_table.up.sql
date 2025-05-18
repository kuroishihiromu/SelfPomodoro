-- ラウンドテーブルの作成
CREATE TABLE IF NOT EXISTS rounds (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    round_order INTEGER NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    work_time INTEGER,
    break_time INTEGER,
    focus_score INTEGER,
    is_skipped BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- インデックスの作成
CREATE INDEX IF NOT EXISTS idx_rounds_session_id ON rounds(session_id);
CREATE INDEX IF NOT EXISTS idx_rounds_session_id_round_order ON rounds(session_id, round_order);

-- トリガーの設定
CREATE TRIGGER rounds_updated_at_trigger
BEFORE UPDATE ON rounds
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
