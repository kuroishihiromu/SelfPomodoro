from typing import Union, List, Optional
import boto3

class DynamoDBHandler:
    """DynamoDBを操作するハンドラークラス"""
    
    def __init__(self, table_name: str, region_name: str) -> None:
        """DynamoDBハンドラーの初期化
        
        Args:
            table_name: 操作するDynamoDBテーブル名
        """
        self.client = boto3.client('dynamodb', region_name=region_name)
        self.table_name = table_name

    def put_round_data(self, user_id: str, time: str, work_time: Optional[float] = None, break_time: Optional[float] = None, focus_score: Optional[float] = None, timestamp: Optional[str] = None) -> dict:
        """ラウンドデータをDynamoDBに追加"""
        try:
            # 必須フィールド
            item = {
                'user_id': {'S': user_id},
                'time': {'S': time}
            }
            
            # オプショナルフィールドをNoneでない場合のみ追加
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

    def get_round_data(self, user_id: str, time: Optional[str] = None) -> Optional[dict] | List[dict]:
        """ラウンドデータを取得"""
        try:
            if time:
                # 特定のアイテムを取得
                response = self.client.get_item(
                    TableName=self.table_name,
                    Key={
                        'user_id': {'S': user_id},
                        'time': {'S': time}
                    }
                )
                return response.get('Item')
            else:
                # ユーザーの全ラウンドデータを取得
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


    def put_session_data(self, user_id: str, time: str, round_count: Optional[int] = None, break_time: Optional[float] = None, avg_focus_score: Optional[float] = None, total_work_time: Optional[float] = None, timestamp: Optional[str] = None) -> dict:
        """セッションデータをDynamoDBに追加"""
        try:
            # 必須フィールド
            item = {
                'user_id': {'S': user_id},
                'time': {'S': time}
            }
            
            # オプショナルフィールドをNoneでない場合のみ追加
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


    def get_session_data(self, user_id: str, time: Optional[str] = None) -> Optional[dict] | List[dict]:
        """セッションデータを取得"""
        try:
            if time:
                # 特定のアイテムを取得
                response = self.client.get_item(
                    TableName=self.table_name,
                    Key={
                        'user_id': {'S': user_id},
                        'time': {'S': time}
                    }
                )
                return response.get('Item')
            else:
                # ユーザーの全セッションデータを取得
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


    def _convert_dynamodb_item_to_dict(self, item: dict) -> dict:
        """DynamoDBのアイテムを通常の辞書形式に変換"""
        if not item:
            return {}
        
        converted = {}
        for key, value in item.items():
            if 'S' in value:  # 文字列
                converted[key] = value['S']
            elif 'N' in value:  # 数値
                # 整数か浮動小数点数かを判定
                if '.' in value['N']:
                    converted[key] = float(value['N'])
                else:
                    converted[key] = int(value['N'])
            # 他のデータ型が必要な場合はここに追加
        return converted


    def _convert_dynamodb_items_to_list(self, items: List[dict]) -> List[dict]:
        """DynamoDBのアイテムリストを通常の辞書のリストに変換"""
        return [self._convert_dynamodb_item_to_dict(item) for item in items]


    def make_chosen_data_list(self, user_id: str, columns: List[str], time: Optional[str] = None) -> Union[List[Union[float, int]], List[List[Union[float, int]]]]:
        """DynamoDBから特定のカラムのデータをリスト形式で取得
        
        Args:
            user_id: ユーザーID
            columns: 取得したいカラム名のリスト
            time: 特定の時間を指定する場合
        
        Returns:
            指定したカラムのデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
        
        Usage:
            columns=["work_time", "break_time"]  # [[25, 5], [30, 5], [20, 5]]
            columns=["work_time"]                # [25, 30, 20]
        """
        try:
            # データ取得（ラウンドデータかセッションデータかは使用時に決定）
            if time:
                # 特定のアイテムを取得
                raw_data = self.get_round_data(user_id, time)
                if not raw_data:
                    return []
                
                converted_data = self._convert_dynamodb_item_to_dict(raw_data)
                # 指定されたカラムの値を取得
                result = [converted_data.get(col) for col in columns if col in converted_data]
                
                # 一つの列の場合は1次元リストとして返す
                if len(columns) == 1 and result:
                    return result[0] if result else []
                return result
            else:
                # 全データを取得
                raw_items = self.get_round_data(user_id)
                if not raw_items:
                    return []
                
                converted_items = self._convert_dynamodb_items_to_list(raw_items)
                
                # 指定されたカラムの値を取得
                result = []
                for item in converted_items:
                    row_data = [item.get(col) for col in columns if col in item]
                    if row_data:
                        result.append(row_data)
                
                # 一つの列の場合は1次元リストにフラット化
                if len(columns) == 1:
                    result = [item[0] for item in result if item]
                
                return result
                
        except Exception as e:
            print(f"データのリスト化に失敗しました: {e}")
            raise Exception(f"データのリスト化に失敗しました: {e}")


    def make_chosen_session_data_list(self, user_id: str, columns: List[str], time: Optional[str] = None) -> Union[List[Union[float, int]], List[List[Union[float, int]]]]:
        """DynamoDBからセッションデータの特定のカラムをリスト形式で取得
        
        Args:
            user_id: ユーザーID
            columns: 取得したいカラム名のリスト
            time: 特定の時間を指定する場合
        
        Returns:
            指定したカラムのデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
        """
        try:
            # データ取得
            if time:
                # 特定のアイテムを取得
                raw_data = self.get_session_data(user_id, time)
                if not raw_data:
                    return []
                
                converted_data = self._convert_dynamodb_item_to_dict(raw_data)
                # 指定されたカラムの値を取得
                result = [converted_data.get(col) for col in columns if col in converted_data]
                
                # 一つの列の場合は1次元リストとして返す
                if len(columns) == 1 and result:
                    return result[0] if result else []
                return result
            else:
                # 全データを取得
                raw_items = self.get_session_data(user_id)
                if not raw_items:
                    return []
                
                converted_items = self._convert_dynamodb_items_to_list(raw_items)
                
                # 指定されたカラムの値を取得
                result = []
                for item in converted_items:
                    row_data = [item.get(col) for col in columns if col in item]
                    if row_data:
                        result.append(row_data)
                
                # 一つの列の場合は1次元リストにフラット化
                if len(columns) == 1:
                    result = [item[0] for item in result if item]
                
                return result
                
        except Exception as e:
            print(f"セッションデータのリスト化に失敗しました: {e}")
            raise Exception(f"セッションデータのリスト化に失敗しました: {e}") 
