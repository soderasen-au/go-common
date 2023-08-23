package exec

import (
	"github.com/soderasen-au/go-common/loggers"
	"github.com/soderasen-au/go-common/util"
)

type InMemMetaKeeper struct {
	metas map[string]Meta
}

func NewInMemMetaKeeper() *InMemMetaKeeper {
	k := &InMemMetaKeeper{}
	k.metas = make(map[string]Meta)
	return k
}

func (k *InMemMetaKeeper) Get(reqId string) (Meta, bool) {
	r, ok := k.metas[reqId]
	return r, ok
}

func (k *InMemMetaKeeper) Set(m Meta) {
	k.metas[m.ID()] = m
}

type InMemRequestKeeper struct {
	keeper MetaKeeper
}

func NewInMemRequestKeeper() *InMemRequestKeeper {
	k := new(InMemRequestKeeper)
	k.keeper = NewInMemMetaKeeper()
	return k
}

func (k *InMemRequestKeeper) Register(req Request) (*RequestMeta, *util.Result) {
	if reqMeta, ok := k.keeper.Get(req.ID()); ok {
		if reqMeta.GetStatus() == StautsRunning {
			return nil, util.MsgError("UpdateMeta", "request is still running")
		}
	}

	reqMeta := &RequestMeta{}
	reqMeta.Reset(req)
	k.keeper.Set(reqMeta)
	return reqMeta, nil
}

func (k *InMemRequestKeeper) GetMeta(id string) (Meta, bool) {
	return k.keeper.Get(id)
}

// caller should run in a co-routine
func (k *InMemRequestKeeper) AsyncRun(req Request) {
	reqLogger := req.Logger()
	if reqLogger == nil {
		reqLogger = loggers.NullLogger
	}
	logger := reqLogger.With().Str("AsyncRun", req.ID()).Logger()

	logger.Info().Msg("start")
	if meta, ok := k.keeper.Get(req.ID()); ok {
		if meta.GetStatus() != StatusReady {
			logger.Error().Msg("request is not ready to run. did you register it first? this request will be IGNORED!")
			return
		}

		meta.SetStatus(StautsRunning)
		succeeded, results := req.Run()
		if !succeeded {
			meta.SetStatus(StatusFailed)
		} else {
			meta.SetStatus(StatusOk)
		}
		meta.SetResults(results)
	} else {
		logger.Error().Msg("can't find request meta. did you register it first? this request will be IGNORED!")
	}
}
