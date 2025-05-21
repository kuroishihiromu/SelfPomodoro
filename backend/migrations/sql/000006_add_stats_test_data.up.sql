-- テストデータを追加するためのマイグレーションスクリプト
-- 1. セッションのテストデータを作成
-- 2. ラウンドのテストデータを作成（集中度スコア付き）

-- テスト用ユーザーID
DO $$
DECLARE
    test_user_id UUID := '00000000-0000-0000-0000-000000000001';
    
    -- セッションID変数を定義
    session1_id UUID := 'a1111111-1111-1111-1111-111111111111';
    session2_id UUID := 'a2222222-2222-2222-2222-222222222222';
    session3_id UUID := 'a3333333-3333-3333-3333-333333333333';
    session4_id UUID := 'a4444444-4444-4444-4444-444444444444';
    session5_id UUID := 'a5555555-5555-5555-5555-555555555555';
    session6_id UUID := 'a6666666-6666-6666-6666-666666666666';
    session7_id UUID := 'a7777777-7777-7777-7777-777777777777';
    session8_id UUID := 'a8888888-8888-8888-8888-888888888888';
    session9_id UUID := 'a9999999-9999-9999-9999-999999999999';
    session10_id UUID := 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
    
    -- ラウンドID変数
    round_id UUID;
    
    -- 現在時刻をベースにした日付変数
    now_time TIMESTAMP := NOW();
    base_date DATE := (now_time - INTERVAL '14 days')::DATE;
    session_date TIMESTAMP;
    round_date TIMESTAMP;
    
    -- 集中度スコアとラウンド時間
    focus_score INTEGER;
    work_minutes INTEGER := 25;
    break_minutes INTEGER := 5;
    
    -- 統計のためのパラメータ
    session_count INTEGER := 0;
    total_rounds INTEGER := 0;
    avg_focus FLOAT := 0;
    total_work_min INTEGER := 0;
    
BEGIN
    -- 1日目の朝のセッション
    session_date := base_date + INTERVAL '9 hours';
    session_count := 3; -- 3ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 78.333; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session1_id, test_user_id, session_date, session_date + INTERVAL '75 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンドデータの挿入
    -- ラウンド1
    round_date := session_date;
    round_id := 'b1111111-1111-1111-1111-111111111111';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session1_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド2
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b1111111-1111-1111-1111-111111111112';
    focus_score := 70;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session1_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド3
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b1111111-1111-1111-1111-111111111113';
    focus_score := 80;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session1_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 1日目の夕方のセッション
    session_date := base_date + INTERVAL '17 hours';
    session_count := 4; -- 4ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 87.5; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session2_id, test_user_id, session_date, session_date + INTERVAL '2 hours', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンドデータの挿入
    -- ラウンド1
    round_date := session_date;
    round_id := 'b2222222-2222-2222-2222-222222222221';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session2_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド2
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b2222222-2222-2222-2222-222222222222';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session2_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド3
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b2222222-2222-2222-2222-222222222223';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session2_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド4
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b2222222-2222-2222-2222-222222222224';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session2_id, 4, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 3日目の午前中のセッション
    session_date := base_date + INTERVAL '2 days' + INTERVAL '10 hours';
    session_count := 2; -- 2ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 65.0; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session3_id, test_user_id, session_date, session_date + INTERVAL '50 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド1 - 集中度低め
    round_date := session_date;
    round_id := 'b3333333-3333-3333-3333-333333333331';
    focus_score := 60;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session3_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- ラウンド2 - 集中度低め
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b3333333-3333-3333-3333-333333333332';
    focus_score := 70;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session3_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 5日目の午後のセッション
    session_date := base_date + INTERVAL '4 days' + INTERVAL '14 hours';
    session_count := 3; -- 3ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 90.0; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session4_id, test_user_id, session_date, session_date + INTERVAL '75 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 集中度高めのラウンド
    round_date := session_date;
    round_id := 'b4444444-4444-4444-4444-444444444441';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session4_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b4444444-4444-4444-4444-444444444442';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session4_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b4444444-4444-4444-4444-444444444443';
    focus_score := 95;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session4_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 7日目の早朝（7時）のセッション
    session_date := base_date + INTERVAL '6 days' + INTERVAL '7 hours';
    session_count := 2; -- 2ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 97.5; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session5_id, test_user_id, session_date, session_date + INTERVAL '50 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 集中度とても高いラウンド（早朝）
    round_date := session_date;
    round_id := 'b5555555-5555-5555-5555-555555555551';
    focus_score := 95;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session5_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b5555555-5555-5555-5555-555555555552';
    focus_score := 100;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session5_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 8日目の深夜（23時）のセッション
    session_date := base_date + INTERVAL '7 days' + INTERVAL '23 hours';
    session_count := 2; -- 2ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 72.5; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session6_id, test_user_id, session_date, session_date + INTERVAL '50 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 深夜の集中度（やや低め）
    round_date := session_date;
    round_id := 'b6666666-6666-6666-6666-666666666661';
    focus_score := 75;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session6_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b6666666-6666-6666-6666-666666666662';
    focus_score := 70;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session6_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 10日目の午後（15時）のセッション
    session_date := base_date + INTERVAL '9 days' + INTERVAL '15 hours';
    session_count := 4; -- 4ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 83.75; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session7_id, test_user_id, session_date, session_date + INTERVAL '2 hours', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 午後のセッション（波のある集中度）
    round_date := session_date;
    round_id := 'b7777777-7777-7777-7777-777777777771';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session7_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b7777777-7777-7777-7777-777777777772';
    focus_score := 75;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session7_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b7777777-7777-7777-7777-777777777773';
    focus_score := 80;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session7_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b7777777-7777-7777-7777-777777777774';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session7_id, 4, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 12日目の午前（11時）のセッション
    session_date := base_date + INTERVAL '11 days' + INTERVAL '11 hours';
    session_count := 3; -- 3ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 85.0; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session8_id, test_user_id, session_date, session_date + INTERVAL '75 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 午前のセッション（安定した集中度）
    round_date := session_date;
    round_id := 'b8888888-8888-8888-8888-888888888881';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session8_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b8888888-8888-8888-8888-888888888882';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session8_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b8888888-8888-8888-8888-888888888883';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session8_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 13日目の夕方（18時）のセッション
    session_date := base_date + INTERVAL '12 days' + INTERVAL '18 hours';
    session_count := 3; -- 3ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 76.67; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session9_id, test_user_id, session_date, session_date + INTERVAL '75 minutes', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 夕方のセッション（徐々に低下する集中度）
    round_date := session_date;
    round_id := 'b9999999-9999-9999-9999-999999999991';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session9_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b9999999-9999-9999-9999-999999999992';
    focus_score := 75;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session9_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'b9999999-9999-9999-9999-999999999993';
    focus_score := 70;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session9_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;

    -- 本日の午前（10時）のセッション（最新データ）
    session_date := now_time - INTERVAL '5 hours';
    session_count := 4; -- 4ラウンド
    total_work_min := session_count * work_minutes;
    avg_focus := 88.75; -- 平均集中度

    -- セッションデータの挿入
    INSERT INTO sessions (id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at)
    VALUES (session10_id, test_user_id, session_date, session_date + INTERVAL '2 hours', avg_focus, total_work_min, session_count, break_minutes * (session_count - 1), session_date, session_date)
    ON CONFLICT (id) DO NOTHING;
    
    -- 本日の午前のセッション（高い集中度）
    round_date := session_date;
    round_id := 'baaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
    focus_score := 90;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session10_id, 1, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'baaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session10_id, 2, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'baaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaac';
    focus_score := 95;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session10_id, 3, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
    round_date := round_date + INTERVAL '30 minutes';
    round_id := 'baaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaad';
    focus_score := 85;
    INSERT INTO rounds (id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at)
    VALUES (round_id, session10_id, 4, round_date, round_date + INTERVAL '25 minutes', work_minutes, break_minutes, focus_score, FALSE, round_date, round_date)
    ON CONFLICT (id) DO NOTHING;
    
END
$$;
