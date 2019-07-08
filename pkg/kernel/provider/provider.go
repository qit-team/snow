package provider

type Provider interface {
	Register(args ... interface{}) (error)
	Provides() ([]string)
	Close() (error)
}
