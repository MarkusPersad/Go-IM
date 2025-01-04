package zaplog

import (
	"fmt"
	"github.com/cilium/lumberjack/v2"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"strings"
	"time"
)

var Logger *zap.Logger

func init() {
	Logger = NewLogger()
}
func NewLogger() *zap.Logger {
	maxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	maxAge, _ := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
	maxBackups, _ := strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
	isCompress, _ := strconv.ParseBool(os.Getenv("LOG_COMPRESS"))
	lv := os.Getenv("LOG_LEVEL")
	currentDate := time.Now().Format("2006-01-02")               // 格式化为 YYYY-MM-DD
	logFileName := fmt.Sprintf("./logs/log-%s.log", currentDate) // 例如 log-2025-01-04.log
	hook := lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   isCompress,
	}
	var level zapcore.Level
	switch strings.ToLower(lv) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	var encoder zapcore.Encoder
	if strings.ToLower(os.Getenv("LOG_ENCODING")) == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	var core zapcore.Core
	if level == zap.DebugLevel {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), level)
	} else {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), level)
	}
	return zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(os.Getenv("LOG_TIME_FORMAT")))
}
