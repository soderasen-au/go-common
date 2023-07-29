package exec

import (
	"github.com/Click-CI/common/loggers"
	"github.com/Click-CI/common/util"
)

type RequestKeeper interface {
	Get(reqId string) (*RequestMeta, bool)
	Set(req *RequestMeta)
}

type BookKeeper struct {
	Requests RequestKeeper
}

func NewInMemBookKeeper() (*BookKeeper, *util.Result) {
	k := new(BookKeeper)
	k.Requests = NewInMemRequestKeeper()
	return k, nil
}

func (k *BookKeeper) Register(req Request) (*RequestMeta, *util.Result) {
	if reqMeta, ok := k.Requests.Get(req.ID()); ok {
		if reqMeta.Status == REQUEST_STATUS_RUNNING {
			return nil, util.MsgError("UpdateMeta", "request is still running")
		}
	}

	reqMeta := &RequestMeta{}
	reqMeta.Reset(req)
	k.Requests.Set(reqMeta)
	return reqMeta, nil
}

func (k *BookKeeper) AsyncRun(req Request) {
	reqLogger := req.Logger()
	if reqLogger == nil {
		reqLogger = loggers.NullLogger
	}
	logger := reqLogger.With().Str("AsyncRun", req.ID()).Logger()

	logger.Info().Msg("start")
	if reqMeta, ok := k.Requests.Get(req.ID()); ok {
		if reqMeta.Status != REQUEST_STATUS_READY {
			logger.Error().Msg("request is not ready to run. did you register it first? this request will be IGNORED!")
			return
		}

		reqMeta.Status = REQUEST_STATUS_RUNNING
		succeeded, results := req.Run()
		if !succeeded {
			reqMeta.Status = REQUEST_STATUS_FAILED
		} else {
			reqMeta.Status = REQUEST_STATUS_OK
		}
		reqMeta.Results = results
	} else {
		logger.Error().Msg("can't find request meta. did you register it first? this request will be IGNORED!")
	}
}
