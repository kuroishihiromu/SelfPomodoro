from skopt.space import Integer
from skopt.learning import GaussianProcessRegressor
from skopt.learning.gaussian_process.kernels import Matern
from skopt import Optimizer

class BayesianOptimizer:
    def __init__(self, target: str):
        """最適化の初期化

        Args:
            target (str): 最適化タイプ
        """
        try:
            if target == "round":
                self.opt = Optimizer(
                    dimensions=[(15.0, 60.0), (3.0, 20.0)],
                    base_estimator=GaussianProcessRegressor(kernel=Matern(length_scale=1.0)),
                    n_initial_points=10,
                    acq_func="EI"
                )
            elif target == "session":
                self.opt = Optimizer(
                    dimensions=[(60.0, 480.0), (10.0, 60.0), Integer(1, 10)],
                    base_estimator=GaussianProcessRegressor(kernel=Matern(length_scale=1.0)),
                    n_initial_points=10,
                    acq_func="EI"
                )
            else:
                raise ValueError(f"最適化タイプが不正です->\n '{target}'")
        except Exception as e:
            print(f"最適化の初期化に失敗しました->\n {e}")
            raise Exception(f"最適化の初期化に失敗しました->\n {e}")
    
    
    def optimize_round(self, explanatory_variable: list[float], objective_variable: list[float]):
        """ラウンド最適化

        Args:
            explanatory_variable (list[float]): 説明変数
            objective_variable (list[float]): 目的変数
        
        Returns:
            tuple[float, float]: 最適化結果
        """
        try:
            # 目的変数を負の値に変換
            negative_objective = [-v for v in objective_variable]
            # データを追加
            self.opt.tell(explanatory_variable, negative_objective)
            # 最適化
            work_time, break_time = self.opt.ask()

            return work_time, break_time
        
        except Exception as e:
            print(f"ラウンド最適化に失敗しました->\n {e}")
            raise Exception(f"ラウンド最適化に失敗しました->\n {e}")
    
    def optimize_session(self, explanatory_variable: list[float], objective_variable: list[float]):
        """セッション最適化

        Args:
            explanatory_variable (list[float]): 説明変数
            objective_variable (list[float]): 目的変数

        Returns:
            tuple[float, float, int]: 最適化結果
        """
        try:
            # 目的変数を負の値に変換
            negative_objective = [-v for v in objective_variable]
            # データを追加
            self.opt.tell(explanatory_variable, negative_objective)
            # 最適化
            total_work_time, break_time, round_count = self.opt.ask()

            return total_work_time, break_time, round_count

        except Exception as e:
            print(f"セッション最適化に失敗しました->\n {e}")
            raise Exception(f"セッション最適化に失敗しました->\n {e}")
