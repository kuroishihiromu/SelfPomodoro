from handler.csv.csv_handler import CSVHandler


def test_csv_handler():
    """CSVHandlerのテスト"""

    csv_handler = CSVHandler("../data/round/123e4567-e89b-12d3-a456-426614174000.csv")

    # --- 説明変数と目的変数を更新 ---
    csv_handler.update_data([90], ["focus_score"])
    csv_handler.update_data([20, 10], ["work_time", "break_time"])

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["work_time", "break_time"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["focus_score"])

    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)

    # --- 型チェック ---
    assert isinstance(explanatory_variable, list)
    assert isinstance(objective_variable, list)

    # --- 値チェック ---
    assert explanatory_variable[-1] == [20, 10]
    assert objective_variable[-2] == 90

    # --- テストで追加したデータを削除 ---
    csv_handler.df = csv_handler.df.iloc[:-1]
    csv_handler.df.iloc[-1, csv_handler.df.columns.get_loc('focus_score')] = None
    csv_handler.df.to_csv(csv_handler.file_path, index=False)
    