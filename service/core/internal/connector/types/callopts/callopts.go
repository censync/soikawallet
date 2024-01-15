package callopts

type CallOpts struct {
	method string
	params []interface{}
}

func NewCallOpts(method string, params []interface{}) *CallOpts {
	return &CallOpts{method: method, params: params}
}

func (c CallOpts) Method() string {
	return c.method
}

func (c CallOpts) Params() []interface{} {
	return c.params
}
