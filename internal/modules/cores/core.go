package cores

type Core interface {
	Initialize(rawConfig map[string]interface{}) error
	Run(runtimeErrorReceiver func(error))
	Stop() error
	GetType() string
}
