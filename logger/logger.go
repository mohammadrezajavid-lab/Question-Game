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

type Config struct {
	FileName               string        `mapstructure:"file_name"`
	MaxSize                int           `mapstructure:"max_size"`
	MaxAge                 int           `mapstructure:"max_age"`
	MaxBackups             int           `mapstructure:"max_backups"`
	Compress               bool          `mapstructure:"compress"`
	SimplingCoreTick       time.Duration `mapstructure:"simpling_core_tick"`
	SimplingCoreFirst      int           `mapstructure:"simpling_core_first"`
	SimplingCoreThereafter int           `mapstructure:"simpling_core_thereafter"`
	Level                  string        `mapstructure:"level"` // e.g. "debug", "info", "warn", "error"
}

func InitLogger(cfg Config) {

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  false,
		Compress:   cfg.Compress,
	})

	var level zapcore.Level
	uErr := level.UnmarshalText([]byte(cfg.Level))
	if uErr != nil {
		zap.L().Fatal(uErr.Error())
	}
	levelAtom := zap.NewAtomicLevelAt(level)
	//level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, levelAtom),
		zapcore.NewCore(fileEncoder, file, levelAtom),
	)

	samplingCore := zapcore.NewSamplerWithOptions(core, cfg.SimplingCoreTick, cfg.SimplingCoreFirst, cfg.SimplingCoreThereafter)

	logger := zap.New(samplingCore)
	zap.ReplaceGlobals(logger)
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

func Warn(err error, msg string, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("warning", err.Error()),
	}

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, args[i+1]))
	}

	zap.L().Named(GetPackageFuncName(2)).Warn(msg, fields...)
}

func Panic(err error, msg string, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("panic", err.Error()),
	}

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, args[i+1]))
	}

	zap.L().Named(GetPackageFuncName(2)).Panic(msg, fields...)
}

func Fatal(err error, msg string, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("fatal", err.Error()),
	}

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, args[i+1]))
	}

	zap.L().Named(GetPackageFuncName(2)).Fatal(msg, fields...)
}

func Info(msg string, args ...interface{}) {
	fields := make([]zapcore.Field, 0)

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, args[i+1]))
	}

	zap.L().Named(GetPackageFuncName(2)).Info(msg, fields...)
}

func Error(err error, msg string, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("error", err.Error()),
	}

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, args[i+1]))
	}

	zap.L().Named(GetPackageFuncName(2)).Error(msg, fields...)
}

func Close() {
	_ = zap.L().Sync()
}
