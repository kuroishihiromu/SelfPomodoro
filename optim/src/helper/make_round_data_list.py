from typing import List, Union


def make_round_data_list(
    round_data: List[dict],
    columns: List[str]
) -> Union[List[Union[float, int]], List[List[Union[float, int]]]]:
    """ラウンドデータから特定のカラムをリスト形式で取得
    
    get_round_data_listと同じ構造のデータを返す関数
    
    Parameters:
        round_data (List[dict]): ラウンドデータのリスト
        columns (List[str]): 取得したいカラム名のリスト
        
    Returns:
        指定したカラムのデータ。一つの列の場合は1次元リスト、複数列の場合は2次元リスト
    """
    if not round_data:
        return []
    
    # 指定されたカラムの値を取得
    result = []
    for item in round_data:
        row_data = [item.get(col) for col in columns if col in item]
        if row_data:
            result.append(row_data)
    
    # 1列の場合は1次元リストにフラット化
    if len(columns) == 1:
        result = [item[0] for item in result if item]
    
    return result
