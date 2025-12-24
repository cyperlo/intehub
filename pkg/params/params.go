package params

type Params map[string]interface{}

func NewParams() Params {
	return make(Params)
}
