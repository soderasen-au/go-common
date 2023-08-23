package exec

import (
	"github.com/soderasen-au/go-common/util"
)

type RequestMeta struct {
	RequestID string         `json:"request_id,omitempty" yaml:"request_id,omitempty" bson:"request_id,omitempty"`
	FuncName  string         `json:"func_name,omitempty" yaml:"func_name,omitempty" bson:"func_name,omitempty"`
	Results   []*util.Result `json:"results,omitempty" yaml:"result,omitempty" bson:"results,omitempty"`
	Status    Status         `json:"status,omitempty" yaml:"status,omitempty" bson:"status,omitempty"`
}

func (m *RequestMeta) GetResults() []*util.Result {
	return m.Results
}

func (m *RequestMeta) SetResults(r []*util.Result) {
	m.Results = r
}

func (m *RequestMeta) SetStatus(s Status) {
	m.Status = s
}

func (m *RequestMeta) ID() string {
	return m.RequestID
}

func (m *RequestMeta) Name() string {
	return m.FuncName
}

func (m *RequestMeta) GetStatus() Status {
	return m.Status
}

func (meta *RequestMeta) Reset(req Request) {
	meta.RequestID = req.ID()
	meta.FuncName = req.Name()
	meta.Results = nil
	meta.Status = StatusReady
}
