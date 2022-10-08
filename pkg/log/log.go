package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"undersea/pkg/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger = zerolog.Logger{}

type ctxLogKeyType string

var ctxLogKey ctxLogKeyType = "logFields"

func fromCtxLogItems(ctx context.Context) map[string]string {
	raw := ctx.Value(ctxLogKey)
	if raw == nil {
		return map[string]string{}
	}
	return raw.(map[string]string)
}

func appendEvents(e *zerolog.Event, addCaller bool) *zerolog.Event {
	e.Timestamp()
	if addCaller {
		_, file, line := funcFileLine("dev.rcrai.com/rcrai/", "dev.rcrai.com/rcrai/zeus-rtareport/core")
		e.Str("caller", fmt.Sprintf("%s:%d", file, line))
	}
	return e
}

func funcFileLine(include string, exclude string) (string, string, int) {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	ff := runtime.CallersFrames(pcs[:n])

	var fn, file string
	var line int
	for {
		f, ok := ff.Next()
		if !ok {
			break
		}
		fn, file, line = f.Function, f.File, f.Line
		if strings.Contains(fn, include) && !strings.Contains(fn, exclude) {
			break
		}
	}

	if ind := strings.LastIndexByte(fn, '/'); ind != -1 {
		fn = fn[ind+1:]
	}

	return fn, file, line
}

var level2Str map[zerolog.Level]string = map[zerolog.Level]string{
	zerolog.DebugLevel: "DEBUG",
	zerolog.InfoLevel:  "INFO",
	zerolog.WarnLevel:  "WARN",
	zerolog.ErrorLevel: "ERROR",
	zerolog.FatalLevel: "FATAL",
	zerolog.PanicLevel: "PANIC",
	zerolog.NoLevel:    "NO",
	zerolog.Disabled:   "DISABLED",
	zerolog.TraceLevel: "TRACE",
}

type logMetricHook struct{}

func initLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	// 统一的把日志时间调整到前面, 所以这儿不用在加了
	Logger = zerolog.New(zerolog.MultiLevelWriter(os.Stdout)).
		With().
		Logger().
		Level(zerolog.InfoLevel)

	level := viper.V().GetString("log.level")
	switch strings.ToUpper(level) {
	case "TRACE":
		Logger = Logger.Level(zerolog.TraceLevel)
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "DEBUG":
		Logger = Logger.Level(zerolog.DebugLevel)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "WARN":
		Logger = Logger.Level(zerolog.WarnLevel)
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		Logger = Logger.Level(zerolog.ErrorLevel)
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		Logger = Logger.Level(zerolog.InfoLevel)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = Logger
}

var once sync.Once
func InitLogger() {
	once.Do(initLogger)
}