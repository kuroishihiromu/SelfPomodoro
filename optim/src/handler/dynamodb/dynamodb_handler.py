from datetime import datetime
from typing import Union, List, Optional
import boto3

class DynamoDBHandler:
    """DynamoDBを操作するハンドラークラス"""

    def __init__(
        self,
        table_name: str,
        region_name: str
    ) -> None:
        """DynamoDBクライアントを作成
        
        Parameters:
            table_name (str): 操作するDynamoDBテーブル名
            region_name (str): AWSリージョン
        """
        self.client = boto3.client('dynamodb', region_name=region_name)
        self.table_name = table_name


    def _convert_to_dict(self, item: dict) -> dict:
        """辞書型に変換"""
        if not item:
            return {}
        
        converted = {}
        for key, value in item.items():
            if 'S' in value:
                converted[key] = value['S']
            elif 'N' in value:
                if '.' in value['N']:
                    converted[key] = float(value['N'])
                else:
                    converted[key] = int(value['N'])
        return converted


    def _convert_to_list(self, items: List[dict]) -> List[dict]:
        """リスト型に変換"""
        return [self._convert_to_dict(item) for item in items]


    def put_round_data(
        self,
        user_id: str,
        time: str,
        work_time: Optional[float] = None,
        break_time: Optional[float] = None,
        focus_score: Optional[float] = None,
        timestamp: Optional[str] = None
    ) -> dict:
        """ラウンドデータをDynamoDBに追加
        
        Parameters:
            user_id (str): ユーザーID
            time (str): 時間
            work_time (Optional[float]): 作業時間
            break_time (Optional[float]): 休憩時間
            focus_score (Optional[float]): 集中度
            timestamp (Optional[str]): タイムスタンプ
        
        Returns:
            dict: 追加したデータ
        """
        try:
            # --- 必須フィールド ---
            item = {
                'user_id': {'S': user_id},
                'time': {'S': time}
            }
            
            # --- オプショナルフィールド ---
            if work_time is not None:
                item['work_time'] = {'N': str(work_time)}
            if break_time is not None:
                item['break_time'] = {'N': str(break_time)}
            if focus_score is not None:
                item['focus_score'] = {'N': str(focus_score)}
            if timestamp is not None:
                item['timestamp'] = {'S': timestamp}
            
            response = self.client.put_item(
                TableName=self.table_name,
                Item=item
            )
            return response
            
        except Exception as e:
            print(f"ラウンドデータの追加に失敗しました: {e}")
            raise Exception(f"ラウンドデータの追加に失敗しました: {e}")


    def put_session_data(
        self,
        user_id: str,
        time: str,
        round_count: Optional[int] = None,
        break_time: Optional[float] = None,
        avg_focus_score: Optional[float] = None,
        total_work_time: Optional[float] = None,
        timestamp: Optional[str] = None
    ) -> dict:
        """セッションデータをDynamoDBに追加
        
        Parameters:
            user_id (str): ユーザーID
            time (str): 時間
            round_count (Optional[int]): ラウンド数
            break_time (Optional[float]): 休憩時間
            avg_focus_score (Optional[float]): 平均集中度
            total_work_time (Optional[float]): 総作業時間
            timestamp (Optional[str]): タイムスタンプ
        
        Returns:
            dict: 追加したデータ
        """
        try:
            # --- 必須フィールド ---
            item = {
                'user_id': {'S': user_id},
                'time': {'S': time}
            }
            
            # --- オプショナルフィールド ---
            if round_count is not None:
                item['round_count'] = {'N': str(round_count)}
            if break_time is not None:
                item['break_time'] = {'N': str(break_time)}
            if avg_focus_score is not None:
                item['avg_focus_score'] = {'N': str(avg_focus_score)}
            if total_work_time is not None:
                item['total_work_time'] = {'N': str(total_work_time)}
            if timestamp is not None:
                item['timestamp'] = {'S': timestamp}
            
            response = self.client.put_item(
                TableName=self.table_name,
                Item=item
            )
            return response
            
        except Exception as e:
            print(f"セッションデータの追加に失敗しました: {e}")
            raise Exception(f"セッションデータの追加に失敗しました: {e}")


    def get_round_data(
        self,
        user_id: str,
        time: Optional[str] = None
    ) -> Optional[dict] | List[dict]:
        """ラウンドデータを取得
        
        Parameters:
            user_id (str): ユーザーID
            time (Optional[str]): 特定の時間を指定する場合
        
        Returns:
            Optional[dict] | List[dict]: ラウンドデータ
        """
        try:
            # --- 特定のアイテムを取得 ---
            if time:
                response = self.client.get_item(
                    TableName=self.table_name,
                    Key={
                        'user_id': {'S': user_id},
                        'time': {'S': time}
                    }
                )
                return response.get('Item')
            
            # --- 全ラウンドデータを取得 ---
            else:
                response = self.client.query(
                    TableName=self.table_name,
                    KeyConditionExpression='user_id = :user_id',
                    ExpressionAttributeValues={
                        ':user_id': {'S': user_id}
                    }
                )
                return response.get('Items', [])
                
        except Exception as e:
            print(f"ラウンドデータの取得に失敗しました: {e}")
            raise Exception(f"ラウンドデータの取得に失敗しました: {e}")


    def get_session_data(
        self,
        user_id: str,
        time: Optional[str] = None
    ) -> Optional[dict] | List[dict]:
        """セッションデータを取得
        
        Parameters:
            user_id (str): ユーザーID
            time (Optional[str]): 特定の時間を指定する場合
        
        Returns:
            Optional[dict] | List[dict]: セッションデータ
        """
        try:
            # --- 特定のアイテムを取得 ---
            if time:
                response = self.client.get_item(
                    TableName=self.table_name,
                    Key={
                        'user_id': {'S': user_id},
                        'time': {'S': time}
                    }
                )
                return response.get('Item')
            
            # --- 全セッションデータを取得 ---
            else:
                response = self.client.query(
                    TableName=self.table_name,
                    KeyConditionExpression='user_id = :user_id',
                    ExpressionAttributeValues={
                        ':user_id': {'S': user_id}
                    }
                )
                return response.get('Items', [])
                
        except Exception as e:
            print(f"セッションデータの取得に失敗しました: {e}")
            raise Exception(f"セッションデータの取得に失敗しました: {e}")


    def get_round_data_list(
        self,
        user_id: str,
        columns: List[str],
        time: Optional[str] = None
    ) -> Union[List[Union[float, int]], List[List[Union[float, int]]]]:
        """特定のカラムデータをリストで取得
        
        Parameters:
            user_id (str): ユーザーID
            columns (List[str]): 取得したいカラム名のリスト
            time (Optional[str]): 特定の時間を指定する場合
        
        Returns:
            指定したカラムのデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
        
        Usage:
            columns=["work_time", "break_time"]
            columns=["work_time"]
        """
        try:
            # --- 特定のアイテムを取得 ---
            if time:
                raw_data = self.get_round_data(user_id, time)
                if not raw_data:
                    return []
                # 辞書型に変換
                converted_data = self._convert_to_dict(raw_data)
                # 指定されたカラムの値を取得
                result = [converted_data.get(col) for col in columns if col in converted_data]
                # 1列の場合は1次元リストにフラット化
                if len(columns) == 1 and result:
                    return result[0] if result else []
                
                return result
            
            # --- 全データを取得 ---
            else:
                raw_items = self.get_round_data(user_id)
                if not raw_items:
                    return []
                # リスト型に変換
                converted_items = self._convert_to_list(raw_items)
                # 指定されたカラムの値を取得
                result = []
                for item in converted_items:
                    row_data = [item.get(col) for col in columns if col in item]
                    if row_data:
                        result.append(row_data)
                # 1列の場合は1次元リストにフラット化
                if len(columns) == 1:
                    result = [item[0] for item in result if item]
                
                return result
                
        except Exception as e:
            print(f"データのリスト化に失敗しました: {e}")
            raise Exception(f"データのリスト化に失敗しました: {e}")


    def get_session_data_list(
        self,
        user_id: str,
        columns: List[str],
        time: Optional[str] = None
    ) -> Union[List[Union[float, int]], List[List[Union[float, int]]]]:
        """セッションデータの特定のカラムをリストで取得
        
        Parameters:
            user_id (str): ユーザーID
            columns (List[str]): 取得したいカラム名のリスト
            time (Optional[str]): 特定の時間を指定する場合
        
        Returns:
            指定したカラムのデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
        """
        try:
            # --- 特定のアイテムを取得 ---
            if time:
                raw_data = self.get_session_data(user_id, time)
                if not raw_data:
                    return []
                # 辞書型に変換
                converted_data = self._convert_to_dict(raw_data)
                # 指定されたカラムの値を取得
                result = [converted_data.get(col) for col in columns if col in converted_data]
                # 1列の場合は1次元リストにフラット化
                if len(columns) == 1 and result:
                    return result[0] if result else []
                return result
            
            # --- 全データを取得 ---
            else:
                raw_items = self.get_session_data(user_id)
                if not raw_items:
                    return []
                # リスト型に変換
                converted_items = self._convert_to_list(raw_items)
                # 指定されたカラムの値を取得
                result = []
                for item in converted_items:
                    row_data = [item.get(col) for col in columns if col in item]
                    if row_data:
                        result.append(row_data)
                # 1列の場合は1次元リストにフラット化
                if len(columns) == 1:
                    result = [item[0] for item in result if item]
                
                return result
                
        except Exception as e:
            print(f"セッションデータのリスト化に失敗しました: {e}")
            raise Exception(f"セッションデータのリスト化に失敗しました: {e}") 


    def get_all_timestamps(
        self,
        user_id: str
    ) -> List[str]:
        """指定されたユーザーのすべてのタイムスタンプを取得
        
        Parameters:
            user_id (str): ユーザーID
        
        Returns:
            List[str]: タイムスタンプのリスト
        """
        try:
            # 全データを取得
            response = self.client.query(
                TableName=self.table_name,
                KeyConditionExpression='user_id = :user_id',
                ExpressionAttributeValues={
                    ':user_id': {'S': user_id}
                }
            )
            
            items = response.get('Items', [])
            
            # タイムスタンプのみを抽出
            timestamps = []
            for item in items:
                if 'time' in item and 'S' in item['time']:
                    timestamps.append(item['time']['S'])
            
            return timestamps
            
        except Exception as e:
            print(f"タイムスタンプの取得に失敗しました: {e}")
            raise Exception(f"タイムスタンプの取得に失敗しました: {e}")


    def get_four_past_days_data(
        self,
        user_id: str
    ) -> tuple[List[dict], List[dict], List[dict], List[dict]]:
        """過去4日間のデータを取得
        
        Parameters:
            user_id (str): ユーザーID
        
        Returns:
            tuple[List[dict], List[dict], List[dict]]: (最新のデータ, 最新から１つ古いデータ, 最新から２つ古いデータ, 最新から３つ古いデータ)
        """
        try:
            # 全タイムスタンプを取得
            all_timestamps = self.get_all_timestamps(user_id)
            
            if not all_timestamps:
                return [], [], [], []
            
            # タイムスタンプを日付でグループ化
            date_groups = {}
            for timestamp in all_timestamps:
                # タイムスタンプから日付部分を抽出
                try:
                    # ISO形式のタイムスタンプをパース
                    if '+' in timestamp:
                        # UTCオフセット付きの場合
                        dt = datetime.fromisoformat(timestamp.replace('+00:00', '+00:00'))
                    else:
                        # オフセットなしの場合
                        dt = datetime.fromisoformat(timestamp)
                    
                    date_key = dt.strftime('%Y-%m-%d')
                    
                    if date_key not in date_groups:
                        date_groups[date_key] = []
                    date_groups[date_key].append(timestamp)
                    
                except ValueError as e:
                    print(f"タイムスタンプのパースに失敗しました: {timestamp}, エラー: {e}")
                    continue
            
            # 過去4日間を取得
            sorted_dates = sorted(date_groups.keys(), reverse=True)
            latest_four_dates = sorted_dates[:4]
            
            # 日付ごとにデータを取得
            latest_data = []
            second_latest_data = []
            third_latest_data = []
            fourth_latest_data = [] 
            
            for i, date in enumerate(latest_four_dates):
                date_data = []
                for timestamp in date_groups[date]:
                    # 各タイムスタンプのデータを取得
                    data = self.get_round_data(user_id, timestamp)
                    if data:
                        date_data.append(data)
                
                if i == 0:
                    latest_data = date_data        # 最新の日付のデータ
                    # ↑ focus_scoreがない最新データ一件（予測用）も含まれる
                elif i == 1:
                    second_latest_data = date_data # 最新から１つ古い日付のデータ
                elif i == 2:
                    third_latest_data = date_data  # 最新から２つ古い日付のデータ
                elif i == 3:
                    fourth_latest_data = date_data # 最新から３つ古い日付のデータ
            
            return latest_data, second_latest_data, third_latest_data, fourth_latest_data
            
        except Exception as e:
            print(f"過去4日間のデータの取得に失敗しました: {e}")
            raise Exception(f"過去4日間のデータの取得に失敗しました: {e}")
