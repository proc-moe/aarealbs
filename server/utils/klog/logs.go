package klog

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

var Log *logType

type logType struct {
	l *logrus.Logger
}

func Init() *logType {
	log := logrus.New()
	ret := &logType{
		l: log,
	}
	Log = ret
	return ret
}

func CtxE(ctx context.Context, format string, a ...any) {
	id := ctx.Value("logId")
	a = append([]any{id}, a...)
	str := fmt.Sprintf("LogId:%v  "+format, a...)
	Log.l.Errorf(str)
}

func CtxW(ctx context.Context, format string, a ...any) {
	id := ctx.Value("logId")
	a = append([]any{id}, a...)
	str := fmt.Sprintf("LogId:%v  "+format, a...)
	Log.l.Warnf(str)
}

func CtxI(ctx context.Context, format string, a ...any) {
	id := ctx.Value("logId")
	a = append([]any{id}, a...)
	str := fmt.Sprintf("LogId:%v  "+format, a...)
	Log.l.Infof(str)
}

func E(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	Log.l.Error(str)
}

func W(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	Log.l.Warnf(str)
}
func I(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	Log.l.Infof(str)
}
