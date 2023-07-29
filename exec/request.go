package exec

import (
	"github.com/Click-CI/common/util"
	"github.com/rs/zerolog"
)

const (
	REQUEST_STATUS_READY   string = "Ready"
	REQUEST_STATUS_FAILED  string = "Failed"
	REQUEST_STATUS_OK      string = "Succeeded"
	REQUEST_STATUS_RUNNING string = "Running"
)

type Request interface {
	ID() string
	Name() string
	Logger() *zerolog.Logger
	Run() (bool, []*util.Result)
}

type RequestMeta struct {
	RequestID string         `json:"request_id,omitempty" yaml:"request_id,omitempty" bson:"request_id,omitempty"`
	FuncName  string         `json:"func_name,omitempty" yaml:"func_name,omitempty" bson:"func_name,omitempty"`
	Results   []*util.Result `json:"results,omitempty" yaml:"result,omitempty" bson:"results,omitempty"`
	Status    string         `json:"status,omitempty" yaml:"status,omitempty" bson:"status,omitempty"`
}

func (meta *RequestMeta) Reset(req Request) {
	meta.RequestID = req.ID()
	meta.FuncName = req.Name()
	meta.Results = nil
	meta.Status = REQUEST_STATUS_READY
}
