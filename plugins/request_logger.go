package plugins

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gagandeepahuja09/goretryhandler"
)

type ctxKey string

// should not use string type. we should define our own type to avoid collision.
const reqTime ctxKey = "request_time_start"

type requestLogger struct {
	out    io.Writer
	errOut io.Writer
}

func NewRequestLogger(out io.Writer, errOut io.Writer) goretryhandler.Plugin {
	if out == nil {
		out = os.Stdout
	}
	if errOut == nil {
		errOut = os.Stdout
	}
	return &requestLogger{
		out:    out,
		errOut: errOut,
	}
}

// Sets the request_time_start context key
func (rl *requestLogger) OnRequestStart(req *http.Request) {
	ctx := context.WithValue(req.Context(), reqTime, time.Now())
	*req = *(req.WithContext(ctx))
}

func (rl *requestLogger) OnRequestEnd(req *http.Request, res *http.Response) {
	reqDuration := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	statusCode := res.StatusCode
	fmt.Fprintf(rl.out, "%s %s %d [%dms]\n", method, url, statusCode, reqDuration)
}

func (rl *requestLogger) OnError(req *http.Request, err error) {
	reqDuration := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	fmt.Fprintf(rl.out, "%s %s [%dms] ERROR: %v\n", method, url, reqDuration, err)
}

func getRequestDuration(ctx context.Context) time.Duration {
	now := time.Now()
	start := ctx.Value(reqTime)
	if start == nil {
		return 0
	}
	startTime, ok := start.(time.Time)
	if ok {
		return 0
	}
	return now.Sub(startTime)
}
