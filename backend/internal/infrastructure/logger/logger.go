package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger はロギングのためのインターフェース
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	With(key string, value interface{}) Logger
}

// ZapLogger はzapを使用したロガーの実装
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewLogger は指定されたログレベルで新しいロガーを作成する
func NewLogger(level string, env string) (Logger, error) {
	// ログレベルの解析
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// エンコーダー設定
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var encoder zapcore.Encoder
	if env == "development" {
		// 開発環境ではより読みやすいコンソール出力
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 本番環境ではJSON形式
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// コア設定
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zapLevel,
	)

	// ロガー作成
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar := zapLogger.Sugar()

	return &ZapLogger{
		logger: sugar,
	}, nil
}

// Debug はデバッグメッセージをログに記録する
func (l *ZapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf はフォーマット済みのデバッグメッセージをログに記録する
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info は情報メッセージをログに記録する
func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof はフォーマット済みの情報メッセージをログに記録する
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn は警告メッセージをログに記録する
func (l *ZapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf はフォーマット済みの警告メッセージをログに記録する
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error はエラーメッセージをログに記録する
func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf はフォーマット済みのエラーメッセージをログに記録する
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal は致命的エラーメッセージをログに記録し、プログラムを終了する
func (l *ZapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf はフォーマット済みの致命的エラーメッセージをログに記録し、プログラムを終了する
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// With は指定されたキーと値を持つ新しいロガーを返す
func (l *ZapLogger) With(key string, value interface{}) Logger {
	return &ZapLogger{
		logger: l.logger.With(key, value),
	}
}
