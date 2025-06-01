from app.services.csv.csv_handler import CSVHandler

def test_round_csv_handler():
    try:
        print("[テスト] CSVファイルの読み込みと更新テスト")
        print()
        # CSVファイルの読み込み
        round_csv_path = "app/data/round/123e4567-e89b-12d3-a456-426614174000.csv"
        csv_handler = CSVHandler(round_csv_path)
        print("CSVファイルの読み込み:\n", csv_handler.df.tail(3))
        print()

        # データの更新
        new_data = [66]
        columns = ["focus_score"]
        csv_handler.update_data(new_data=new_data, columns=columns)
        print("CSVファイルの更新:\n", csv_handler.df.tail(3))
        print()

        # 説明変数と目的変数を取得
        explanatory_variable = csv_handler.make_chosen_data_list(columns=["work_time", "break_time"])
        objective_variable = csv_handler.make_chosen_data_list(columns=["focus_score"])
        print("説明変数リストの取得:\n", explanatory_variable)
        print()
        print("目的変数リストの取得:\n", objective_variable)
        assert isinstance(explanatory_variable, list)
        assert isinstance(objective_variable, list)
        print()
        print()
        print("[成功]: CSVファイルの読み込みと更新テスト完了")
        print()
        
    except FileNotFoundError as e:
        print(f"[失敗]: CSVファイルが見つかりません: {e}")
        print()
        raise
    except ValueError as e:
        print(f"[失敗]: データの更新に失敗しました: {e}")
        print()
        raise
    except Exception as e:
        print(f"[失敗]: 予期しないエラーが発生しました: {e}")
        print()
        raise

if __name__ == "__main__":
    test_round_csv_handler()