import os
import pandas as pd

class CSVHandler:
    """CSVファイルを操作するクラス

    Attributes:
        file_path (str): CSVファイルの位置
        df(DataFrame): CSVファイルのデータ
    Methods:
        __init__(file_path: str)
         -> CSVファイルの読み込み
        update_data(new_data: list, columns: list[str] = None)
         -> CSVファイルの更新
        choose_data(columns: list[str] = None, rows: list[int] | slice = None)
         -> CSVファイルの特定のカラムと行を取得
        make_chosen_data_list(columns: list[str] = None, rows: list[int] | slice = None)
         -> CSVファイルの特定のカラムと行を取得し、リストに変換
    """
    def __init__(self, file_path):
        """ファイルの読み込み

        Args:
            file_path (str): CSVファイルの位置
        """
        try:
            # CSVファイルのパスを定義
            self.file_path = file_path
            # CSVファイルの存在確認
            if not os.path.exists(file_path):
                raise FileNotFoundError(f"File not found: {file_path}")
            # CSVファイルの読み込み
            self.df = pd.read_csv(self.file_path)

        except Exception as e:
            print(f"CSVファイルの読み込みに失敗しました->\n {e}")
            raise {f"CSVファイルの読み込みに失敗しました->\n {e}"}


    def update_data(self, new_data: list, columns: list[str] = None):
        """CSVファイルの更新

        Args:
            new_data (list): 新しいデータ
            columns (list[str]): 更新するカラム名
        """
        try:
            # columnsが文字列（単一カラム）の場合、リストに変換
            if isinstance(columns, str):
                columns = [columns]
                new_data = [new_data]
            # 最後の行の状態を確認
            last_row = self.df.iloc[-1]
            has_empty_columns = last_row.isna().any()
            if has_empty_columns and all(col in self.df.columns for col in columns):
                # 最後の行に空の列があり、新しいデータがその列に対応する場合
                # 既存の最後の行にデータを追加
                last_index = len(self.df) - 1
                for col, val in zip(columns, new_data):
                    self.df.at[last_index, col] = val
            else:
                # そうでないなら新しい行として追加
                new_row = pd.DataFrame([new_data], columns=columns)
                self.df = pd.concat([self.df, new_row], ignore_index=True)
            # CSVファイルの更新
            self.df.to_csv(self.file_path, index=False)
            
        except Exception as e:
            print(f"CSVファイルの更新に失敗しました->\n {e}")
            raise Exception(f"CSVファイルの更新に失敗しました->\n {e}")


    def choose_data(self, columns: list[str] = None, rows: list[int] | slice = None):
        """CSVファイルの特定のカラムと行を取得

        Args:
            columns (list[str]): 取得したいカラム
            rows (list[int] | slice): 取得したい行
        Useage:
            columns=["work_time", "break_time"], rows=[0,1,4]
            columns=["work_time", "break_time"], rows=slice(0,2)
            columns=["work_time", "break_time"], rows=slice(44,None)
            columns=["work_time"], rows=None
            columns=None, rows=slice(0,2)
        Returns:
            chosen_data (DataFrame): 特定のカラムと行のデータ
        """
        try:
            # カラムと行の取得
            if columns and rows is not None:
                chosen_data = self.df.iloc[rows][columns]
            # カラムのみの取得            
            elif columns is not None and rows is None:
                chosen_data = self.df[columns]
            # 行のみの取得
            elif rows is not None and columns is None:
                chosen_data = self.df.iloc[rows]
            
            else:
                raise ValueError("カラムか列のいずれかを指定してください")
            
            return chosen_data
        
        except Exception as e:
            print(f"CSVファイルの特定のカラムと行の取得に失敗しました->\n {e}")
            raise {f"CSVファイルの特定のカラムと行の取得に失敗しました->\n {e}"}
    

    def make_chosen_data_list(self, columns: list[str] = None, rows: list[int] | slice = None):
        """CSVファイルの特定のカラムと行を取得し、リストに変換

        Args:
            columns (list[str]): 取得したいカラム
            rows (list[int] | slice): 取得したい行
        Useage:
            columns=["work_time", "break_time"], rows=[0,1,4]  // [[25, 5], [30, 5], [20, 5]]
            columns=["work_time"], rows=None                   // [25, 30, 20]
            columns=["work_time", "break_time"], rows=slice(0,2)
            columns=["work_time", "break_time"], rows=slice(44,None)
            columns=None, rows=slice(0,2)
        Returns:
            chosen_data_list (list): 特定のカラムと行のデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
        """
        try:
            # カラムと行の取得
            chosen_data = self.choose_data(columns=columns, rows=rows)
            # リストに変換
            chosen_data_list = chosen_data.values.tolist()
            # 一つの列だけの場合は1次元リストにフラット化
            if columns is not None and len(columns) == 1:
                chosen_data_list = [item[0] for item in chosen_data_list]
            
            return chosen_data_list
        
        except Exception as e:
            print(f"CSVファイルの特定のカラムと行の取得に失敗しました->\n {e}")
            raise {f"CSVファイルの特定のカラムと行の取得に失敗しました->\n {e}"}
    
    
if __name__ == "__main__":
    # CSVファイル操作
    csv_handler = CSVHandler("../../data/round_data.csv")
    # CSVファイルの更新
    new_data = [6]
    columns = ["focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)
    # 説明変数と目的変数を取得
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["work_time", "break_time"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)