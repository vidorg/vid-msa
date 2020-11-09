package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	LogFileDir    string //文件保存地方
	AppName       string //日志文件前缀
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	Level         zapcore.Level //日志等级
	MaxSize       int           //日志文件小大（M）
	MaxBackups    int           // 最多存在多少个切片文件
	MaxAge        int           //保存的最大天数
	Development   bool          //是否是开发模式
	zap.Config
}

type ModOptions func(options *Options)

var (
	Log                            *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type Logger struct {
	*zap.Logger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func NewLogger(mod ...ModOptions) *zap.Logger {
	Log = &Logger{}
	Log.Lock()
	defer Log.Unlock()
	if Log.inited {
		Log.Info("[Logger] logger Inited")
		return nil
	}
	Log.Opts = &Options{
		LogFileDir:    "",
		AppName:       "app_log",
		ErrorFileName: "error.log",
		WarnFileName:  "warn.log",
		InfoFileName:  "info.log",
		DebugFileName: "debug.log",
		Level:         zapcore.DebugLevel,
		MaxSize:       100,
		MaxBackups:    60,
		MaxAge:        30,
	}
	if Log.Opts.LogFileDir == "" {
		Log.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		Log.Opts.LogFileDir += sp + "log" + sp
	}
	if Log.Opts.Development {
		Log.zapConfig = zap.NewDevelopmentConfig()
		Log.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		Log.zapConfig = zap.NewProductionConfig()
		Log.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	if Log.Opts.OutputPaths == nil || len(Log.Opts.OutputPaths) == 0 {
		Log.zapConfig.OutputPaths = []string{"stdout"}
	}
	if Log.Opts.ErrorOutputPaths == nil || len(Log.Opts.ErrorOutputPaths) == 0 {
		Log.zapConfig.OutputPaths = []string{"stderr"}
	}
	for _, fn := range mod {
		fn(Log.Opts)
	}
	Log.zapConfig.Level.SetLevel(Log.Opts.Level)
	Log.init()
	Log.inited = true
	Log.Info("[Logger] success")
	return Log.Logger
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}
func SetMaxSize(MaxSize int) ModOptions {
	return func(option *Options) {
		option.MaxSize = MaxSize
	}
}
func SetMaxBackups(MaxBackups int) ModOptions {
	return func(option *Options) {
		option.MaxBackups = MaxBackups
	}
}
func SetMaxAge(MaxAge int) ModOptions {
	return func(option *Options) {
		option.MaxAge = MaxAge
	}
}

func SetLogFileDir(LogFileDir string) ModOptions {
	return func(option *Options) {
		option.LogFileDir = LogFileDir
	}
}

func SetAppName(AppName string) ModOptions {
	return func(option *Options) {
		option.AppName = AppName
	}
}

func SetLevel(Level zapcore.Level) ModOptions {
	return func(option *Options) {
		option.Level = Level
	}
}
func SetErrorFileName(ErrorFileName string) ModOptions {
	return func(option *Options) {
		option.ErrorFileName = ErrorFileName
	}
}
func SetWarnFileName(WarnFileName string) ModOptions {
	return func(option *Options) {
		option.WarnFileName = WarnFileName
	}
}

func SetInfoFileName(InfoFileName string) ModOptions {
	return func(option *Options) {
		option.InfoFileName = InfoFileName
	}
}
func SetDebugFileName(DebugFileName string) ModOptions {
	return func(option *Options) {
		option.DebugFileName = DebugFileName
	}
}
func SetDevelopment(Development bool) ModOptions {
	return func(option *Options) {
		option.Development = Development
	}
}
func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	//consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}

// Info Info
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn Warn
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal Fatal
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// Panic Panic
func Panic(msg string, fields ...zap.Field) {
	Log.Panic(msg, fields...)
}
