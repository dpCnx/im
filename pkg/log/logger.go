package log

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"im/api/pb"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(c *pb.Log) *Logger {
	return &Logger{log: newZapLogger(c)}
}

func newZapLogger(c *pb.Log) *zap.Logger {
	infoWriteSyncer := initInfoLogWriter(c)
	encoder := getEncoder()
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.DebugLevel
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriteSyncer, infoPriority),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)) // 根据上面的配置创建logger
	return logger
	// zap.ReplaceGlobals(logger)               // 替换zap库里全局的logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig) // json格式日志
}

func initInfoLogWriter(c *pb.Log) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.FileName,
		MaxSize:    int(c.MaxSize),    // 日志文件大小 单位：MB
		MaxBackups: int(c.MaxBackups), // 备份数量
		MaxAge:     int(c.MaxAge),     // 备份时间 单位：天
		Compress:   true,              // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

/*
根据时间切割	假设 rotationTime 为 24hour 即「天」，那么当你的程序在任意时间启动，都会在半夜 12 点，产生日志切割。
*/
func getLogWriter() zapcore.WriteSyncer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.logger.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotateLogs.New(
		fmt.Sprintf("./logger/%s", "demo.%Y_%m_%d.info"),
		rotateLogs.WithMaxAge(time.Duration(1)*time.Hour*24),
		rotateLogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(hook)
}

func (l *Logger) Log(level log.Level, keyVal ...interface{}) error {
	args := l.argsToFile(keyVal)

	switch level {
	case log.LevelDebug:
		l.log.Debug("", args...)
	case log.LevelInfo:
		l.log.Info("", args...)
	case log.LevelWarn:
		l.log.Warn("", args...)
	case log.LevelError:
		l.log.Error("", args...)
	}
	return nil
}

func (l *Logger) argsToFile(args ...interface{}) []zap.Field {
	res := args[0].([]interface{})
	if len(res) <= 0 {
		return nil
	}

	if len(res)%2 != 0 {
		l.log.Error(fmt.Sprint("args must appear in pairs: ", args))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(res); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(res[i]), fmt.Sprint(res[i+1])))
	}
	return data
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, l.argsToFile(args)...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, l.argsToFile(args)...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, l.argsToFile(args)...)
}
