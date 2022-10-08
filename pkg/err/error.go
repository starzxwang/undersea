package error_v2

import (
	"context"
	"encoding/json"
	"undersea/pkg/log"

	"github.com/pborman/uuid"
	"runtime"
	"time"
)

// 为了统一管理error上下文链路信息，这里对error封装一层

// 定义一个Error接口
type ErrorInterface interface {
	// 打印error信息
	getTraceId() string
	error
}

type Error struct {
	traceId string
}

// 每个error处的代码位置信息
type caller struct {
	// trace_id
	TraceId string `json:"trace_id"`
	// 行号
	LineNo int `json:"line_no"`
	// 所在函数名
	FuncName string `json:"func_name"`
	// 错误信息
	Msg string `json:"msg"`
	// 当前文件名
	FileName string `json:"file_name"`
	// 当前时间
	CurrentTime time.Time `json:"-"`
}

func initError() ErrorInterface {
	return &Error{
		traceId: uuid.New(),
	}
}

func (me *Error) getTraceId() string {
	return me.traceId
}

func printCaller(ctx context.Context, call *caller, err error) {
	b, _ := json.Marshal(call)
	log.E(ctx, err).Msgf(call.CurrentTime.String()+" "+string(b))
}

//func PrintError(ctx context.Context, err error)  {
//	for _, v := range err.(*Error).callers {
//		printCaller(ctx, v, err)
//	}
//}

func (me *Error) Error() string {
	return ""
}

func PrintError(ctx context.Context, err error, msg string) error {
	if err == nil {
		err = initError()
	}

	// skip=1，即可以拿到调用InjectError函数的代码位置信息
	pc, fullFilePath, lineNo, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	printCaller(ctx, &caller{
		FileName: fullFilePath,
		LineNo: lineNo,
		Msg: msg,
		CurrentTime: time.Now(),
		FuncName: runtime.FuncForPC(pc).Name(),
		TraceId: err.(*Error).getTraceId(),
	}, err)

	return err
}