package forcecast

type ForceCast interface {
	Cast(value interface{}) (interface{}, error)
	ZeroValue() interface{}
}
