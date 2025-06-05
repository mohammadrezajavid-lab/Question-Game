package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    30,
		MaxAge:     7,
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   true,
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	samplingCore := zapcore.NewSamplerWithOptions(core, time.Minute, 100, 100)

	logger := zap.New(samplingCore)
	zap.ReplaceGlobals(logger)
}

func GetPackageName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullFuncName := fn.Name() // github.com/user/project/pkg/module.Func
	parts := strings.Split(fullFuncName, "/")
	last := parts[len(parts)-1] // module.Func
	pkgParts := strings.Split(last, ".")
	if len(pkgParts) > 1 {
		return pkgParts[0]
	}

	return "unknown"
}

func GetPackageFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullFuncName := fn.Name() // github.com/user/project/pkg/module.Func
	parts := strings.Split(fullFuncName, "/")
	last := parts[len(parts)-1] // module.Func

	return last
}
