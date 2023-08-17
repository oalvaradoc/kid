package callback

// Options is defined all the options that callback executor can set
type Options struct {
	Port                      int
	CommType                  string // mandatory(mesh/direct): default mesh
	ServerAddress             string // mandatory: default http://127.0.0.1:18080
	CallbackPort              int    // mandatory: default 18082
	EnableClientSideStatusFSM bool
	ExtConfigs                map[string]interface{}
}

var (
	// DefaultPort is default port of http endpoint listener
	DefaultPort = 6060

	// DefaultCommType is default communication type
	DefaultCommType = "mesh"

	// DefaultServerAddress  is default address of SED server
	DefaultServerAddress = "http://127.0.0.1:18080"

	// DefaultCallbackPort  is default port of callback endpoint listener
	DefaultCallbackPort = 18082

	// DefaultEnableClientSideStatusFSM is default of whether enable client side status FSM.
	DefaultEnableClientSideStatusFSM = false
)

// Option is the type of closure function that defines the registration settings Options
type Option func(*Options)

// NewHandlerOptions is create a callback executor handler with default setting.
func NewHandlerOptions(options ...Option) Options {
	opts := Options{
		Port:                      DefaultPort,
		CommType:                  DefaultCommType,
		ServerAddress:             DefaultServerAddress,
		CallbackPort:              DefaultCallbackPort,
		EnableClientSideStatusFSM: DefaultEnableClientSideStatusFSM,
		ExtConfigs:                make(map[string]interface{}),
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

// WithPort is used to modify the port number
func WithPort(port int) Option {
	return func(options *Options) {
		options.Port = port
	}
}

// WithCommType is used to modify the communicate type
func WithCommType(commType string) Option {
	return func(options *Options) {
		options.CommType = commType
	}
}

// WithServerAddress is used to modify the server address
func WithServerAddress(serverAddress string) Option {
	return func(options *Options) {
		options.ServerAddress = serverAddress
	}
}

// WithCallbackPort is used to modify the callback port
func WithCallbackPort(callbackPort int) Option {
	return func(options *Options) {
		options.CallbackPort = callbackPort
	}
}

// WithEnableClientSideStatusFSM is used to open or close the status FSM of client side
func WithEnableClientSideStatusFSM(enableClientSideStatusFSM bool) Option {
	return func(options *Options) {
		options.EnableClientSideStatusFSM = enableClientSideStatusFSM
	}
}

// WithExtConfigs is used to setting the additional configs.
func WithExtConfigs(extConfigs map[string]interface{}) Option {
	return func(options *Options) {
		options.ExtConfigs = extConfigs
	}
}

// AddExtConfig is used to add the key-value pair config to the exiting additional config.
// This function will create a new config if the additional config is empty
func AddExtConfig(key string, value interface{}) Option {
	return func(options *Options) {
		if nil == options.ExtConfigs {
			options.ExtConfigs = make(map[string]interface{})
		}
		options.ExtConfigs[key] = value
	}
}
