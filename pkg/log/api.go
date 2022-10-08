package log

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func T(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Trace(), false))
}

func D(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Debug(), false))
}

func I(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Info(), false))
}

func W(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Warn(), true))
}

func E(ctx context.Context, err error) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Error(), true)).Err(err)
}

func F(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Fatal(), true))
}

func WithLevel(ctx context.Context, level zerolog.Level) *zerolog.Event {
	addCaller := level > zerolog.InfoLevel
	return WithLogContext(ctx, appendEvents(log.WithLevel(level), addCaller))
}

// WErr
// W(ctx).Err(err).Send()
func WErr(ctx context.Context, err error) {
	W(ctx).Err(err).Send()
}

func GetLevel() zerolog.Level {
	return log.Logger.GetLevel()
}

// WithLogContext
//
// Example: WithLogContext(ctx, log.Warn()).Msgf("hello %s", "huskar")
func WithLogContext(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	if ctx == nil || ctx == context.TODO() || ctx == context.Background() {
		return e
	}

	logFields := fromCtxLogItems(ctx)
	if len(logFields) == 0 {
		return e
	}

	for k, v := range logFields {
		e = e.Str(k, v)
	}

	return e
}

// WithLogValues 通过 context.WithValue 向 ctx 中加入特定 key, value.
// 随后可以通过 WithLogContext 加入到 zerolog.Event 中.
//
// Example: WithLogValues(ctx, "call", "2413", "customer", "2578").
// 当 items 的数量不正确的为奇数时, 忽略最后一个字符串.
func WithLogValues(ctx context.Context, items ...string) context.Context {
	if len(items) == 0 {
		return ctx
	}

	logFields := fromCtxLogItems(ctx)
	for i := 0; i+1 < len(items); i += 2 {
		logFields[items[i]] = items[i+1]
	}

	return context.WithValue(ctx, ctxLogKey, logFields)
}
