import json
import logging
import boto3
import os
from datetime import datetime
from typing import Dict, Any, List, Optional

# Bayesian optimization imports
import numpy as np
from skopt.space import Real, Integer
from skopt.learning import GaussianProcessRegressor
from skopt.learning.gaussian_process.kernels import Matern
from skopt import Optimizer

# ログ設定
logger = logging.getLogger()
logger.setLevel(logging.INFO)

# DynamoDB設定
DYNAMODB_TABLE_NAME = os.environ.get('DYNAMODB_TABLE_NAME', 'selfpomodoro_unified_table_dev')
dynamodb = boto3.client('dynamodb')


class BayesianOptimizer:
    """Bayesian optimization for session parameters using scikit-optimize"""
    
    def __init__(self):
        self.opt = Optimizer(
            dimensions=[
                Real(20.0, 480.0, name='total_work_time'),  # total work time range in minutes (20-480分)
                Real(5.0, 60.0, name='break_time'),         # break time range in minutes (5-60分)
                Integer(1, 10, name='round_count')          # round count range (1-10ラウンド)
            ],
            base_estimator=GaussianProcessRegressor(
                kernel=Matern(length_scale=1.0, nu=2.5),
                alpha=1e-6,
                normalize_y=True
            ),
            n_initial_points=10,
            acq_func="EI",  # Expected Improvement
            acq_func_kwargs={'xi': 0.01}  # Exploration parameter
        )
    
    def optimize_session(self, explanatory_variable: List[List[float]], objective_variable: List[float]) -> tuple:
        """Optimize session parameters using Bayesian optimization"""
        try:
            if not explanatory_variable or not objective_variable:
                logger.info("No historical data available, returning default values")
                return 100.0, 15.0, 4
            
            if len(explanatory_variable) < 3:
                logger.info(f"Insufficient data ({len(explanatory_variable)} points), returning default values")
                return 100.0, 15.0, 4
            
            # Convert to numpy arrays for better performance
            X = np.array(explanatory_variable)
            y = np.array(objective_variable)
            
            logger.info(f"Training Bayesian optimizer with {len(X)} data points")
            logger.info(f"Focus score range: {y.min():.3f} - {y.max():.3f}")
            
            # Add data to optimizer (scikit-optimize expects minimization, so we negate y)
            self.opt.tell(X.tolist(), (-y).tolist())
            
            # Get next suggestion
            next_point = self.opt.ask()
            total_work_time, break_time, round_count = next_point
            
            logger.info(f"Bayesian optimization suggests: total_work_time={total_work_time:.1f}, break_time={break_time:.1f}, round_count={round_count}")
            
            return float(total_work_time), float(break_time), int(round_count)
            
        except Exception as e:
            logger.error(f"Bayesian optimization failed: {str(e)}")
            logger.info("Falling back to default values")
            return 100.0, 15.0, 4


def get_optimization_data(user_id: str) -> tuple:
    """Get historical optimization data from unified DynamoDB table"""
    try:
        # Query optimization log entries for this user
        response = dynamodb.query(
            TableName=DYNAMODB_TABLE_NAME,
            KeyConditionExpression='PK = :pk AND begins_with(SK, :sk_prefix)',
            ExpressionAttributeValues={
                ':pk': {'S': f'USER#{user_id}'},
                ':sk_prefix': {'S': 'OPTIMIZATION_LOG#SESSION#'}
            },
            ScanIndexForward=True  # Sort by SK ascending (oldest first)
        )
        
        items = response.get('Items', [])
        explanatory_vars = []
        objective_vars = []
        
        for item in items:
            # Extract data from DynamoDB item
            total_work_time = item.get('total_work_time', {}).get('N')
            break_time = item.get('break_time', {}).get('N')
            round_count = item.get('round_count', {}).get('N')
            avg_focus_score = item.get('avg_focus_score', {}).get('N')
            
            if total_work_time and break_time and round_count and avg_focus_score:
                explanatory_vars.append([float(total_work_time), float(break_time), int(float(round_count))])
                objective_vars.append(float(avg_focus_score))
        
        return explanatory_vars, objective_vars
        
    except Exception as e:
        logger.error(f"Failed to get optimization data: {str(e)}")
        return [], []


def update_user_config(user_id: str, session_break_time: float, session_rounds: int) -> None:
    """Update OptimizationPreferences with session optimization results"""
    try:
        timestamp = datetime.now().isoformat()
        
        # Update OptimizationPreferences with new optimized values
        dynamodb.put_item(
            TableName=DYNAMODB_TABLE_NAME,
            Item={
                'PK': {'S': f'USER#{user_id}'},
                'SK': {'S': 'OPTIMIZATION_PREFERENCES'},
                'user_id': {'S': user_id},
                'round_work_time': {'N': '25'},  # Default value
                'round_break_time': {'N': '5'},  # Default value
                'session_rounds': {'N': str(session_rounds)},
                'session_break_time': {'N': str(int(session_break_time))},
                'updated_at': {'S': timestamp},
                'entity_type': {'S': 'optimization_preferences'}
            },
            # Update existing config or create if not exists
            ConditionExpression='attribute_exists(PK) AND attribute_exists(SK)'
        )
        
        logger.info(f"Updated OptimizationPreferences for user {user_id[:8]}... - session_break_time: {session_break_time:.1f}, session_rounds: {session_rounds}")
        
    except Exception as e:
        # If preferences doesn't exist, create it
        if 'ConditionalCheckFailedException' in str(e):
            logger.info(f"OptimizationPreferences not found for user {user_id[:8]}..., creating new one")
            try:
                dynamodb.put_item(
                    TableName=DYNAMODB_TABLE_NAME,
                    Item={
                        'PK': {'S': f'USER#{user_id}'},
                        'SK': {'S': 'OPTIMIZATION_PREFERENCES'},
                        'user_id': {'S': user_id},
                        'round_work_time': {'N': '25'},  # Default value
                        'round_break_time': {'N': '5'},  # Default value
                        'session_rounds': {'N': str(session_rounds)},
                        'session_break_time': {'N': str(int(session_break_time))},
                        'created_at': {'S': timestamp},
                        'updated_at': {'S': timestamp},
                        'entity_type': {'S': 'optimization_preferences'}
                    }
                )
                logger.info(f"Created new OptimizationPreferences for user {user_id[:8]}...")
            except Exception as create_error:
                logger.error(f"Failed to create OptimizationPreferences: {str(create_error)}")
                raise
        else:
            logger.error(f"Failed to update OptimizationPreferences: {str(e)}")
            raise


def save_optimization_result(user_id: str, session_id: str, total_work_time: float, break_time: float, round_count: int, avg_focus_score: float) -> None:
    """Save optimization result to unified DynamoDB table"""
    try:
        timestamp = datetime.now().isoformat()
        
        # Save optimization log entry
        dynamodb.put_item(
            TableName=DYNAMODB_TABLE_NAME,
            Item={
                'PK': {'S': f'USER#{user_id}'},
                'SK': {'S': f'OPTIMIZATION_LOG#SESSION#{timestamp}'},
                'entity_type': {'S': 'optimization_log'},
                'optimization_type': {'S': 'session'},
                'total_work_time': {'N': str(total_work_time)},
                'break_time': {'N': str(break_time)},
                'round_count': {'N': str(round_count)},
                'avg_focus_score': {'N': str(avg_focus_score)},
                'created_at': {'S': timestamp},
                'updated_at': {'S': timestamp}
            }
        )
        
        logger.info(f"Saved optimization result for user {user_id[:8]}...")
        
    except Exception as e:
        logger.error(f"Failed to save optimization result: {str(e)}")
        raise


def lambda_handler(event: Dict[str, Any], context) -> Dict[str, Any]:
    """
    SQSトリガー型セッション最適化Lambda関数（Bayesian optimization版）
    
    Args:
        event: SQS event containing session optimization message
        context: Lambda context
    
    Returns:
        Dict containing processing results
    """
    try:
        logger.info(f"セッション最適化処理開始: {len(event.get('Records', []))} messages")
        
        results = []
        optimizer = BayesianOptimizer()
        
        # SQSメッセージを順次処理
        for record in event.get('Records', []):
            try:
                # SQSメッセージからデータを解析
                message_body = json.loads(record['body'])
                logger.info(f"処理対象メッセージ: {message_body.get('message_id', 'unknown')}")
                
                # 必要なデータを抽出
                user_id = message_body['user_id']
                session_id = message_body['session_id']
                avg_focus_score = message_body['avg_focus_score']
                
                logger.info(f"最適化開始 - ユーザー: {user_id[:8]}..., セッション: {session_id[:8]}..., 平均スコア: {avg_focus_score}")
                
                # 過去の最適化データを取得
                explanatory_vars, objective_vars = get_optimization_data(user_id)
                logger.info(f"取得データ: {len(explanatory_vars)} records")
                
                # Bayesian optimization実行
                total_work_time, break_time, round_count = optimizer.optimize_session(explanatory_vars, objective_vars)
                
                logger.info(f"最適化完了 - 作業時間: {total_work_time:.1f}分, 休憩時間: {break_time:.1f}分, ラウンド数: {round_count}")
                
                # 結果をDynamoDBに保存
                save_optimization_result(user_id, session_id, total_work_time, break_time, round_count, avg_focus_score)
                
                # OptimizationPreferencesをセッション最適化結果で更新
                update_user_config(user_id, break_time, round_count)
                
                optimization_result = {
                    'total_work_time': total_work_time,
                    'break_time': break_time,
                    'round_count': round_count
                }
                
                results.append({
                    'message_id': message_body.get('message_id'),
                    'user_id': user_id,
                    'session_id': session_id,
                    'status': 'success',
                    'optimization_result': optimization_result
                })
                
            except Exception as e:
                logger.error(f"メッセージ処理エラー: {str(e)}")
                results.append({
                    'message_id': record.get('messageId', 'unknown'),
                    'status': 'error',
                    'error': str(e)
                })
        
        # 処理結果のサマリー
        success_count = len([r for r in results if r['status'] == 'success'])
        error_count = len([r for r in results if r['status'] == 'error'])
        
        logger.info(f"セッション最適化処理完了 - 成功: {success_count}, エラー: {error_count}")
        
        return {
            'statusCode': 200,
            'body': {
                'message': 'セッション最適化処理完了',
                'processed_count': len(results),
                'success_count': success_count,
                'error_count': error_count,
                'results': results
            }
        }
        
    except Exception as e:
        logger.error(f"Lambda関数実行エラー: {str(e)}")
        return {
            'statusCode': 500,
            'body': {
                'message': 'セッション最適化処理失敗',
                'error': str(e)
            }
        }