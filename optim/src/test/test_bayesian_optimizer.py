from optimizer.bayesian_optimizer import BayesianOptimizer


def test_bayesian_optimizer():
    """BayesianOptimizerのテスト"""

    opt = BayesianOptimizer("round")

    # --- 説明変数と目的変数 ---
    explanatory_variable = [[35.94517823265169, 8.58383252694626], [35.811200506982274, 10.169277459808319], [23.812024080081557, 3.4434741687810186], [22.2531636870093, 14.929248985726153], [25.155074559971656, 14.854696396224353], [20, 10], [20, 10], [30.46009388692101, 18.834463857625806]]
    objective_variable = [88, 99, 55, 66, 77, 90, 90, 77]

    # --- 説明変数と目的変数を取得 ---
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)

    # --- 型チェック ---
    assert isinstance(work_time, float)
    assert isinstance(break_time, float)

    # --- 値チェック ---
    assert work_time == 47.189713658100594
    assert break_time == 11.892873717275231
