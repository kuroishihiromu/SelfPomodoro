package optimization

import (
	"context"
	"fmt"
	"time"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// SampleOptimizationDataGenerationDomainService はサンプル最適化データ生成のドメインサービス
type SampleOptimizationDataGenerationDomainService struct {
	sessionRepo repository.SessionRepository
	logger      logger.Logger
}

// NewSampleOptimizationDataGenerationDomainService は新しいサンプル最適化データ生成サービスを作成する
func NewSampleOptimizationDataGenerationDomainService(
	sessionRepo repository.SessionRepository,
	logger logger.Logger,
) *SampleOptimizationDataGenerationDomainService {
	return &SampleOptimizationDataGenerationDomainService{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// GenerateAndStoreSampleData はサンプル最適化データを生成し、適切なリポジトリに保存する
func (s *SampleOptimizationDataGenerationDomainService) GenerateAndStoreSampleData(ctx context.Context, userID userVO.UserID) error {
	s.logger.Infof("サンプル最適化データ生成開始: UserID=%s", userID.String())

	// 過去10日間のサンプルデータを生成
	baseTime := time.Now().AddDate(0, 0, -10)
	
	// サンプルパターンデータ
	samplePatterns := s.getSamplePatterns()
	
	for _, pattern := range samplePatterns {
		if pattern.sessions == 0 {
			continue // 休日はスキップ
		}
		
		currentDate := baseTime.AddDate(0, 0, pattern.day)
		
		// 1日のセッション数分ループ
		for sessionNum := 0; sessionNum < pattern.sessions; sessionNum++ {
			// セッション開始時刻（9時、13時、17時）
			sessionStartTime := currentDate.Add(time.Duration(9+sessionNum*4) * time.Hour)
			
			// サンプルセッションを作成
			session, err := s.createSampleSession(userID, sessionStartTime, pattern)
			if err != nil {
				s.logger.Errorf("サンプルセッション作成エラー: %v", err)
				continue
			}
			
			// サンプルラウンドを作成・追加・完了
			err = s.createAndCompleteSessionRounds(session, sessionStartTime, pattern)
			if err != nil {
				s.logger.Errorf("サンプルラウンド作成エラー: %v", err)
				continue
			}
			
			// セッションを完了状態に設定（統計情報計算）
			session.CompleteWithRounds()
			
			// リポジトリに保存
			err = s.sessionRepo.CreateSession(ctx, session)
			if err != nil {
				s.logger.Errorf("サンプルセッション保存エラー: %v", err)
				continue
			}
		}
	}
	
	s.logger.Infof("サンプル最適化データ生成完了: UserID=%s", userID.String())
	return nil
}

// createSampleSession はサンプルセッションを作成する
func (s *SampleOptimizationDataGenerationDomainService) createSampleSession(
	userID userVO.UserID,
	startTime time.Time,
	pattern samplePattern,
) (*entity.Session, error) {
	// サンプルセッションの作成
	session := entity.NewSession(userID)
	
	// Note: サンプルデータ用にセッション開始時刻を過去の時刻に設定したい場合、
	// エンティティに適切なメソッドを追加する必要がある
	// 現在は新規作成として扱う
	
	return session, nil
}

// createAndCompleteSessionRounds はセッションにサンプルラウンドを作成・追加・完了する
func (s *SampleOptimizationDataGenerationDomainService) createAndCompleteSessionRounds(
	session *entity.Session,
	sessionStartTime time.Time,
	pattern samplePattern,
) error {
	roundPatterns := s.getRoundPatterns()
	
	for roundNum := 0; roundNum < pattern.roundsPerSession; roundNum++ {
		if roundNum >= len(roundPatterns) {
			break // パターン数を超えた場合は終了
		}
		
		roundPattern := roundPatterns[roundNum]
		
		// セッションにラウンドを追加（maxRoundsは大きな値を設定）
		round, err := session.AddRound(10)
		if err != nil {
			return fmt.Errorf("ラウンド追加エラー: %v", err)
		}
		
		// 作業時間・休憩時間（次回ラウンド用の最適化値）
		workTime := roundPattern.workTime
		breakTime := roundPattern.breakTime
		
		// 経験による微調整（後半の日では最適化が進む）
		if pattern.day > 5 {
			if roundNum == 0 {
				workTime = 30 // 最初のラウンドは30分に延長
			} else if roundNum >= 3 {
				workTime = 20 // 後半は短めに調整
			}
		}
		
		// 集中度スコア（このラウンドの実績値）
		focusScore := int(pattern.baseFocus) + roundPattern.focusAdjust
		
		// 値の範囲チェック
		if focusScore > 100 {
			focusScore = 100
		}
		if focusScore < 30 {
			focusScore = 30
		}
		
		// ラウンドを完了
		err = round.CompleteWith(&focusScore, workTime, breakTime)
		if err != nil {
			return fmt.Errorf("ラウンド完了エラー: %v", err)
		}
	}
	
	return nil
}

// サンプルパターンデータ構造
type samplePattern struct {
	day              int
	sessions         int
	baseFocus        float64
	roundsPerSession int
}

// getSamplePatterns はサンプルパターンデータを返す（ユーザー提供データに基づく）
func (s *SampleOptimizationDataGenerationDomainService) getSamplePatterns() []samplePattern {
	return []samplePattern{
		{0, 2, 72.0, 4},  // 1日目: 2セッション, 集中度72, 4ラウンド
		{1, 1, 68.0, 3},  // 2日目: 1セッション, 集中度68, 3ラウンド
		{2, 3, 75.0, 4},  // 3日目: 3セッション, 集中度75, 4ラウンド
		{3, 2, 70.0, 5},  // 4日目: 2セッション, 集中度70, 5ラウンド
		{4, 1, 65.0, 2},  // 5日目: 1セッション, 集中度65, 2ラウンド
		{5, 2, 78.0, 3},  // 6日目: 2セッション, 集中度78, 3ラウンド
		{6, 0, 0.0, 0},   // 7日目: 休み
		{7, 3, 73.0, 4},  // 8日目: 3セッション, 集中度73, 4ラウンド
		{8, 2, 69.0, 4},  // 9日目: 2セッション, 集中度69, 4ラウンド
		{9, 1, 76.0, 2},  // 10日目: 1セッション, 集中度76, 2ラウンド
		{10, 2, 71.0, 3}, // 11日目: 2セッション, 集中度71, 3ラウンド
		{11, 3, 74.0, 5}, // 12日目: 3セッション, 集中度74, 5ラウンド
		{12, 2, 77.0, 3}, // 13日目: 2セッション, 集中度77, 3ラウンド
	}
}

// ラウンドパターンデータ構造
type roundPattern struct {
	order       int
	workTime    int
	breakTime   int
	focusAdjust int
}

// getRoundPatterns はラウンドパターンデータを返す
func (s *SampleOptimizationDataGenerationDomainService) getRoundPatterns() []roundPattern {
	return []roundPattern{
		{0, 25, 5, +5},   // 1ラウンド目: 集中しやすい
		{1, 25, 5, 0},    // 2ラウンド目: 標準
		{2, 25, 10, -3},  // 3ラウンド目: 少し疲れ
		{3, 20, 15, -5},  // 4ラウンド目: 短めで調整
		{4, 15, 20, -8},  // 5ラウンド目: さらに短めで大幅休憩
		{5, 15, 25, -10}, // 6ラウンド目: 最短で長休憩
	}
}