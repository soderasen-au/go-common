package util

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

func FuncName(skip int) string {
	pc, filePath, _, ok := runtime.Caller(skip)
	if !ok {
		return "InvalidFuncName"
	}
	_, file := filepath.Split(filePath)
	_, funcName := filepath.Split(runtime.FuncForPC(pc).Name())
	return fmt.Sprintf("%s::>%s", file, funcName)
}

type Result struct {
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	Ctx       string      `json:"context"`
	Timestamp time.Time   `json:"timestamp"`
	Result    interface{} `json:"result"`
	Inner     *Result     `json:"inner,omitempty"`
}

func (r *Result) Error() string {
	if r == nil {
		return ""
	}

	str := fmt.Sprintf("%s | %d | %s", r.Ctx, r.Code, r.Msg)
	if r.Inner != nil {
		str = fmt.Sprintf("%s -> [%s]", str, r.Inner.Error())
	}
	return str
}

func Errorf(format string, a ...interface{}) *Result {
	return &Result{
		Code:      -1,
		Msg:       fmt.Sprintf(format, a...),
		Ctx:       FuncName(2),
		Timestamp: time.Now(),
	}
}

func ErrResult(err error) *Result {
	return &Result{
		Code:      -1,
		Msg:       err.Error(),
		Ctx:       FuncName(2),
		Timestamp: time.Now(),
	}
}

func Error(ctx string, err error) *Result {
	return &Result{
		Code:      -1,
		Msg:       err.Error(),
		Ctx:       ctx,
		Timestamp: time.Now(),
	}
}

func MsgError(ctx string, msg string) *Result {
	return &Result{
		Code:      -1,
		Msg:       msg,
		Ctx:       ctx,
		Timestamp: time.Now(),
	}
}

func LogMsgError(logger *zerolog.Logger, ctx string, msg string) *Result {
	logger.Error().Msgf("%s: %s", ctx, msg)
	return &Result{
		Code:      -1,
		Msg:       msg,
		Ctx:       ctx,
		Timestamp: time.Now(),
	}
}

func OK(ctx string) *Result {
	return &Result{
		Code:      0,
		Msg:       "OK",
		Ctx:       ctx,
		Timestamp: time.Now(),
	}
}

func NewResult(ctx string, r interface{}) *Result {
	return &Result{
		Code:      0,
		Msg:       "OK",
		Ctx:       ctx,
		Timestamp: time.Now(),
		Result:    r,
	}
}

func NewErrResult(ctx string, r interface{}) *Result {
	return &Result{
		Code:      -1,
		Msg:       "Failed",
		Ctx:       ctx,
		Timestamp: time.Now(),
		Result:    r,
	}
}

func (r *Result) With(ctx string) *Result {
	return &Result{
		Code:      -1,
		Msg:       "Error",
		Ctx:       ctx,
		Timestamp: time.Now(),
		Inner:     r,
	}
}

func (r *Result) LogWith(logger *zerolog.Logger, ctx string) *Result {
	logger.Error().Err(r).Msg(ctx)
	return &Result{
		Code:      -1,
		Msg:       "Error",
		Ctx:       ctx,
		Timestamp: time.Now(),
		Inner:     r,
	}
}

func (r Result) GetRootCause() *Result {
	if r.Inner != nil {
		return r.Inner.GetRootCause()
	}

	ret := r
	return &ret
}

func Jsonify(r interface{}) json.RawMessage {
	res, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("can't generate response buff: " + err.Error())
	}
	return res
}

func JsonStr(r interface{}) string {
	res, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("can't generate response buff: " + err.Error())
	}
	return string(res)
}

func SameResult(lh, rh *Result) bool {
	if lh == rh {
		return true
	}
	if lh != nil && rh != nil {
		return lh.Code == rh.Code && lh.Ctx == rh.Ctx && lh.Msg == rh.Msg && lh.Result == rh.Result
	}
	return false
}
