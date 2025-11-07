import numpy as np
from sklearn.preprocessing import MinMaxScaler
from tensorflow.keras.models import Sequential # type: ignore
from tensorflow.keras.layers import GRU, Dense # type: ignore

class OldFocusScoreModel:
  """集中度スコアを予測するモデル"""

  def __init__(
      self,
      time_step: int
  ) -> None:
      """GRUモデルの初期化
      Args:
          time_step: 時系列データのステップ数
      """
      # time_stepを保存
      self.time_step = time_step
      
      # スケーラーの初期化
      self.scaler = MinMaxScaler(feature_range=(0, 1))
      
      # GRUモデルの構築
      self.model = Sequential()
      self.model.add(GRU(50, return_sequences=True, input_shape=(self.time_step, 5)))
      self.model.add(GRU(50, return_sequences=False))
      self.model.add(Dense(25))
      self.model.add(Dense(1))

      self.model.compile(optimizer='adam', loss='mean_squared_error')
      
      # 訓練済みフラグ
      self.is_trained = False

  def create_dataset(self, dataset: np.ndarray, time_step: int) -> tuple[np.ndarray, np.ndarray]:
      """時系列データセットの作成"""
      X, Y = [], []
      
      # time_stepが0以下の場合は1に設定
      if time_step <= 0:
          time_step = 1
          print("time_stepが0以下のため、1に設定しました")
      
      # データが少ない場合は、time_stepを調整
      if len(dataset) <= time_step + 1:
          time_step = max(1, len(dataset) - 2)
          print(f"データが少ないため、time_stepを{time_step}に調整しました")
      
      # 最低1つのサンプルが作成されるように調整
      if len(dataset) - time_step - 1 <= 0:
          time_step = max(1, len(dataset) - 2)
          print(f"サンプル作成のため、time_stepを{time_step}に再調整しました")

      for i in range(len(dataset)-time_step-1):
          a = dataset[i:(i+time_step), :]  # 全5次元のデータ
          X.append(a)
          Y.append(dataset[i + time_step, 2])  # focus_score
      
      print(f"データセット作成: データ長={len(dataset)}, time_step={time_step}, 作成されたサンプル数={len(X)}")
      return np.array(X), np.array(Y)

  def fit(
      self,
      train_data: np.ndarray
  ) -> None:
      """モデルの訓練
      Args:
          train_data: 訓練データ
      """
      # スケーリング
      self.scaler.fit(train_data)
      train_data_scaled = self.scaler.transform(train_data)
      
      # データセットの作成
      X_train, y_train = self.create_dataset(train_data_scaled, self.time_step)
      
      print(f"訓練データ形状: X_train={X_train.shape}, y_train={y_train.shape}")
      
      # モデルの訓練
      self.model.fit(X_train, y_train, batch_size=1, epochs=1)
      self.is_trained = True

  def predict(
      self,
      test_data: np.ndarray
  ) -> float:
      """予測
      Args:
          test_data: テストデータ
      Returns:
          test_predict: テストデータに対する予測結果
      """
      if not self.is_trained:
          raise ValueError("モデルが訓練されていません。fit()メソッドを先に実行してください。")
      
      # スケーリング
      test_data_scaled = self.scaler.transform(test_data)
      
      # 最後のtime_step分のデータを使用して予測
      if len(test_data_scaled) >= self.time_step:
          # 最後のtime_step分のデータを取得
          X_test_pred = test_data_scaled[-self.time_step:].reshape(1, self.time_step, 5)
      else:
          # データが少ない場合は、パディングしてtime_step分にする
          padding = np.zeros((self.time_step - len(test_data_scaled), 5))
          X_test_pred = np.vstack([padding, test_data_scaled]).reshape(1, self.time_step, 5)
      
      test_predict = self.model.predict(X_test_pred)

      # スケールを元に戻す
      dummy_data = np.zeros((1, 5))  # スケーリングのため5次元ダミー配列を作成
      dummy_data[0, 2] = test_predict[0, 0]
      test_predict_inverse = self.scaler.inverse_transform(dummy_data)
      predicted_focus_score = test_predict_inverse[0, 2]
      
      return predicted_focus_score
