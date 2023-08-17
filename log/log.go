package log

import (
	"context"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// KeyTraceID is a key for trance ID
	KeyTraceID = "traceId"
	// KeyParentID is a key for parent span ID
	KeyParentID = "parentId"
	// KeySpanID is a key for span ID
	KeySpanID = "spanId"
	// KeyServiceID is a key for service ID
	KeyServiceID = "serviceID"
	// KeyTopicID is a key for topic ID
	KeyTopicID = "topicID"
)

// RotateConfig is used to organize the rotating policy for logger
type RotateConfig struct {
	FilePath   string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// Config is used to organize the config for logger
type Config struct {
	Rotate  RotateConfig
	Level   string
	Console bool
}

// Equals returns whether the self and other are equals
func (c Config) Equals(o *Config) bool {
	return reflect.DeepEqual(&c, o)
}

// Level is a type of logger level
type Level int

const (
	// DebugLevel is used to specify the log level as debug
	DebugLevel = Level(zapcore.DebugLevel)

	// InfoLevel is used to specify the log level as info
	InfoLevel = Level(zapcore.InfoLevel)

	// WarnLevel is used to specify the log level as warn
	WarnLevel = Level(zapcore.WarnLevel)

	// ErrorLevel is used to specify the log level as error
	ErrorLevel = Level(zapcore.ErrorLevel)

	// PanicLevel is used to specify the log level as panic
	PanicLevel = Level(zapcore.PanicLevel)

	// FatalLevel is used to specify the log level as fatal
	FatalLevel = Level(zapcore.FatalLevel)
)

type HookFunc func(ctx context.Context) (context.Context, []interface{})

var (
	logger                    = zap.NewNop()
	sugar                     = logger.Sugar()
	loggerLevel zapcore.Level = zapcore.DebugLevel
	hookFn      HookFunc

	_logLevelMap = map[string]zapcore.Level{
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
		"error": zapcore.ErrorLevel,
		"warn":  zapcore.WarnLevel,
		"info":  zapcore.InfoLevel,
		"debug": zapcore.DebugLevel,
	}
	config = Config{
		Rotate: RotateConfig{
			FilePath:   "/data/logs",
			Filename:   "default.log",
			MaxSize:    200,
			MaxBackups: 7,
			MaxAge:     7,
		},
		Level:   "debug",
		Console: true,
	}
)

func init() {
	var err error
	developmentConfig := zap.NewDevelopmentConfig()
	developmentConfig.EncoderConfig.EncodeCaller = nil
	if logger, err = developmentConfig.Build(); err != nil {
		panic(err)
	}
	sugar = logger.Sugar()
}

func CurrentConfig() Config {
	return config
}

// Init is used to initial the logger
// conf: option parameter, will using the default config if caller without any parameter
func Init(conf ...Config) {
	if len(conf) > 0 {
		config = conf[0]
	}
	rotateConfig := config.Rotate
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	instanceID := os.Getenv("INSTANCE_ID")

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(rotateConfig.FilePath, instanceID, rotateConfig.Filename),
		MaxSize:    rotateConfig.MaxSize,
		MaxBackups: rotateConfig.MaxBackups,
		MaxAge:     rotateConfig.MaxAge,
		Compress:   false,
	}
	fileSync := zapcore.AddSync(lumberJackLogger)
	var (
		cores = make([]zapcore.Core, 1)
	)
	loggerLevel = _logLevelMap[config.Level]
	cores[0] = zapcore.NewCore(encoder, fileSync, loggerLevel)
	if config.Console {
		consoleSync := zapcore.AddSync(os.Stdout)

		cores = append(cores, zapcore.NewCore(encoder, consoleSync, loggerLevel))
	}
	core := zapcore.NewTee(cores...)
	logger = zap.New(core)
	sugar = logger.Sugar()
	defer func() { _ = logger.Sync() }() // flushes buffer, if any
}

// RegisterHookFunc register hook function to add custom log field
func RegisterHookFunc(fn HookFunc) {
	hookFn = fn
}

// IsEnable is used to determine whether the current log level is less than or equal to the specified level
func IsEnable(level Level) bool {
	return level >= Level(loggerLevel)
}

// Debug is used to print log at debug level
// ctx is used to implicit pass the parameters
// msg is the string that need to print
func Debug(ctx context.Context, msg string) {
	if IsEnable(DebugLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Debugw(msg, buildStandardKeyValues(ctx, custom)...)
	}
}

// Info is used to print log at info level
// ctx is used to implicit pass the parameters
// msg is the string that need to print
func Info(ctx context.Context, msg string) {
	if IsEnable(InfoLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Infow(msg, buildStandardKeyValues(ctx, custom)...)
	}
}

// Error is used to print log at error level
// ctx is used to implicit pass the parameters
// msg is the string that need to print
func Error(ctx context.Context, msg string) {
	if IsEnable(ErrorLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Errorw(msg, buildStandardKeyValues(ctx, custom)...)
	}
}

// Warn is used to print log at warn level
// ctx is used to implicit pass the parameters
// msg is the string that need to print
func Warn(ctx context.Context, msg string) {
	if IsEnable(WarnLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Warnw(msg, buildStandardKeyValues(ctx, custom)...)
	}
}

// Debugf is used to print log with format at debug level
// ctx is used to implicit pass the parameters
// format is the string that need to print
// args is an optional parameter with format args
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if IsEnable(DebugLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Debugw(fmt.Sprintf(format, args...), buildStandardKeyValues(ctx, custom)...)
	}
}

// Infof is used to print log with format at info level
// ctx is used to implicit pass the parameters
// format is the string that need to print
// args is an optional parameter with format args
func Infof(ctx context.Context, format string, args ...interface{}) {
	if IsEnable(InfoLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Infow(fmt.Sprintf(format, args...), buildStandardKeyValues(ctx, custom)...)
	}
}

// Errorf is used to print log with format at error level
// ctx is used to implicit pass the parameters
// format is the string that need to print
// args is an optional parameter with format args
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if IsEnable(ErrorLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Errorw(fmt.Sprintf(format, args...), buildStandardKeyValues(ctx, custom)...)
	}
}

// Warnf is used to print log with format at warn level
// ctx is used to implicit pass the parameters
// format is the string that need to print
// args is an optional parameter with format args
func Warnf(ctx context.Context, format string, args ...interface{}) {
	if IsEnable(WarnLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Warnw(fmt.Sprintf(format, args...), buildStandardKeyValues(ctx, custom)...)
	}
}

// Debugs uses fmt.Sprint to construct and log a message at debug level.
func Debugs(args ...interface{}) {
	sugar.Debug(args...)
}

// Infos uses fmt.Sprint to construct and log a message at info level.
func Infos(args ...interface{}) {
	sugar.Info(args...)
}

// Errors uses fmt.Sprint to construct and log a message at error level.
func Errors(args ...interface{}) {
	sugar.Error(args...)
}

// Warns uses fmt.Sprint to construct and log a message at warn level.
func Warns(args ...interface{}) {
	sugar.Warn(args...)
}

// Debugsf uses fmt.Sprintf to log a templated message at debug level.
func Debugsf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

// Infosf uses fmt.Sprintf to log a templated message at info level.
func Infosf(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

// Errorsf uses fmt.Sprintf to log a templated message at error level.
func Errorsf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

// Warnsf uses fmt.Sprintf to log a templated message at warn level.
func Warnsf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// Debugsw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Debugsw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

// Infosw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infosw(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

// Errorsw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorsw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

// Warnsw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnsw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if IsEnable(DebugLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Debugw(msg, append(buildStandardKeyValues(ctx, custom), keysAndValues...)...)
	}
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if IsEnable(InfoLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Infow(msg, append(buildStandardKeyValues(ctx, custom), keysAndValues...)...)
	}
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if IsEnable(ErrorLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Errorw(msg, append(buildStandardKeyValues(ctx, custom), keysAndValues...)...)
	}
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if IsEnable(WarnLevel) {
		var custom []interface{}
		if hookFn != nil {
			ctx, custom = hookFn(ctx)
		}
		sugar.Warnw(msg, append(buildStandardKeyValues(ctx, custom), keysAndValues...)...)
	}
}

type options struct {
	m map[string]string
}

// Option is type of func(o *options)
type Option func(o *options)

// MetaDataWithMap returns an Option which contains an map with key-value pairs
func MetaDataWithMap(m map[string]string) Option {
	return func(o *options) {
		o.m = m
	}
}

// Metadata returns an Option with a key-value pair
func Metadata(key, value string) Option {
	return func(o *options) {
		o.m[key] = value
	}
}

// AuditInfo is used to print audit log at info level
func AuditInfo(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Infow("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}

// AuditError is used to print audit log at error level
func AuditError(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Errorw("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}

// AuditDebug is used to print audit log at debug level
func AuditDebug(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Debugw("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}

func buildAuditKeyValues(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) []interface{} {
	res := buildStandardKeyValues(ctx, nil)
	res = append(res, "user", user, "operator", operator, "request", string(request), "operatorLevel", operatorLevel, "logtype", "audit")
	var o = options{
		m: make(map[string]string),
	}
	for _, f := range metadata {
		f(&o)
	}
	for k, v := range o.m {
		res = append(res, k, v)
	}
	return res
}

func buildStandardKeyValues(ctx context.Context, custom []interface{}) []interface{} {
	var sign int16
	for i := 0; i < len(custom); i += 2 {
		if str, ok := custom[i].(string); ok {
			switch str {
			case KeyTraceID: //1
				sign |= 1
			case KeySpanID: //2
				sign |= (1 << 1)
			case KeyParentID: //3
				sign |= (1 << 2)
			case KeyServiceID: //4
				sign |= (1 << 3)
			case KeyTopicID: //5
				sign |= (1 << 4)
			}
		}
	}
	res := make([]interface{}, 0)
	if traceID := checkNil(ctx.Value(KeyTraceID)); traceID != "" && (sign&1) == 0 {
		res = append(res, KeyTraceID, traceID)
	}
	if spanID := checkNil(ctx.Value(KeySpanID)); spanID != "" && sign&(1<<1) == 0 {
		res = append(res, KeySpanID, spanID)
	}
	if parentID := checkNil(ctx.Value(KeyParentID)); parentID != "" && sign&(1<<2) == 0 {
		res = append(res, KeyParentID, parentID)
	}
	if serviceID := checkNil(ctx.Value(KeyServiceID)); serviceID != "" && sign&(1<<3) == 0 {
		res = append(res, KeyServiceID, serviceID)
	}
	if topicID := checkNil(ctx.Value(KeyTopicID)); topicID != "" && sign&(1<<4) == 0 {
		res = append(res, KeyTopicID, topicID)
	}
	if len(custom) > 0 && len(custom)&1 == 0 {
		res = append(res, custom...)
	}
	return res
}

func checkNil(arg interface{}) interface{} {
	if arg == nil {
		return ""
	}
	return arg
}
