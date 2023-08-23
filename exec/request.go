package exec

import (
	"github.com/rs/zerolog"
	"github.com/soderasen-au/go-common/util"
)

type Status string

const (
	StatusReady   Status = "ready"
	StatusFailed  Status = "failed"
	StatusOk      Status = "succeeded"
	StautsRunning Status = "running"
)

type Request interface {
	ID() string
	Name() string
	Logger() *zerolog.Logger
	Run() (bool, []*util.Result)
}

type Meta interface {
	ID() string
	Name() string
	GetResults() []*util.Result
	SetResults(r []*util.Result)
	GetStatus() Status
	SetStatus(s Status)
}

type MetaKeeper interface {
	Get(reqId string) (Meta, bool)
	Set(req Meta)
}

type RequestKeeper interface {
	Register(req Request) (*RequestMeta, *util.Result)
	GetMeta(id string) (Meta, bool)
	AsyncRun(req Request)
}
