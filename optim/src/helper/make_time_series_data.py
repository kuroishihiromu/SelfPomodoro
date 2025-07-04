from typing import List

def make_time_series_data(
    datas: List[dict],
) -> tuple[List[List[float]], List[List[float]]]:
    """時系列データの作成

    Args:
        datas (List[dict]): データリスト

    Returns:
        tuple[List[List[float]], List[List[float]]]: 時系列データと予測用データ
    """
    try:
        # --- 時系列データの作成 ---
        total_work_time = 0
        total_break_time = 0
        time_series_data = []
        no_focus_score_data = []
        
        for data in datas:
            work_time = data.get('work_time')
            if not work_time:
                raise ValueError("work_timeが存在しません")
            short_break_time = data.get('break_time')
            if not short_break_time:
                raise ValueError("short_break_timeが存在しません")
            work_time = float(work_time.get('N'))               # 作業時間
            short_break_time = float(short_break_time.get('N')) # 短休憩時間
            total_work_time += work_time                        # 総作業時間
            total_break_time += short_break_time                # 総休憩時間
            
            # --- 訓練用と予測用のデータを用意 ---
            focus_score_data = data.get('focus_score')
            if focus_score_data and 'N' in focus_score_data:
                # 訓練用
                focus_score = float(focus_score_data.get('N'))
                time_series_data.append([
                    work_time,
                    short_break_time,
                    focus_score,
                    total_work_time,
                    total_break_time
                ])
            else:
                # 予測用
                focus_score = 0.0  # Noneではなく0.0を使用
                no_focus_score_data.append([
                    work_time,
                    short_break_time,
                    focus_score,
                    total_work_time,
                    total_break_time
                ])

        return time_series_data, no_focus_score_data
    
    except Exception as e:
        print(f"時系列データの作成に失敗しました: {e}")
        raise Exception(f"時系列データの作成に失敗しました: {e}")
