package exec

type InMemRequestKeeper struct {
	requestMap map[string]*RequestMeta
}

func NewInMemRequestKeeper() *InMemRequestKeeper {
	k := &InMemRequestKeeper{}
	k.requestMap = make(map[string]*RequestMeta)
	return k
}

func (k *InMemRequestKeeper) Get(reqId string) (*RequestMeta, bool) {
	r, ok := k.requestMap[reqId]
	return r, ok
}

func (k *InMemRequestKeeper) Set(req *RequestMeta) {
	k.requestMap[req.RequestID] = req
}
