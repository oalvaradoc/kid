package config

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
)

var serviceConfigs *ServiceConfigs

// ConfigOnChangeHookFunc is a controller that used to config on change hook function registration.
type ConfigOnChangeHookFunc struct {
	Name   string
	Func   func(oldConfig *ServiceConfigs, newConfig *ServiceConfigs) error
	IsSync bool
}

var (
	// configOnChangeHookFuncMap defines config on change hook func type variables
	configOnChangeHookFuncMap = make(map[string]*ConfigOnChangeHookFunc)
	locker                    sync.RWMutex
)

// RegisterConfigOnChangeHookFunc provides the registration of a hook function for configuration update events
func RegisterConfigOnChangeHookFunc(name string, handleFunc func(oldConfig *ServiceConfigs, newConfig *ServiceConfigs) error, isSync bool) {
	locker.Lock()
	defer locker.Unlock()
	log.Infosf("Register init hook function, function name[%s], handleFunc[%++v], isSync[%v]", name, handleFunc, isSync)
	configOnChangeHookFuncMap[name] = &ConfigOnChangeHookFunc{
		Name:   name,
		Func:   handleFunc,
		IsSync: isSync,
	}
}

// RunHookHandleFuncIfNecessary is start all the hooked functions when config changed if there is any hook function has registed.
func RunHookHandleFuncIfNecessary(oldConfig *ServiceConfigs, newConfig *ServiceConfigs) error {
	locker.RLock()
	defer locker.RUnlock()
	if len(configOnChangeHookFuncMap) > 0 {
		var wg sync.WaitGroup
		log.Infosf("Service config, start run hook ...")
		for k, v := range configOnChangeHookFuncMap {
			log.Infosf("Service config, start running hook[%s], hook info:[%++v]", k, v)
			if v.IsSync {
				// sync call
				if err := v.Func(oldConfig, newConfig); nil != err {
					log.Errorsf("Service config, running sync hook function[%s] failed, error=%++v", k, err)
					return err
				}
				log.Infosf("Service config, run hook[%s] end!", k)
			} else {
				// async call
				go func(hookName string, configOnChangeHookFunc *ConfigOnChangeHookFunc) {
					defer func() {
						log.Infosf("Service config, run hook[%s] end!", hookName)
						wg.Done()
					}()
					wg.Add(1)
					if err := configOnChangeHookFunc.Func(oldConfig, newConfig); nil != err {
						log.Errorsf("Service config, running async hook function[%s] failed, error=%++v", k, err)
					}
				}(k, v)
			}
		}
		wg.Wait()
		log.Infosf("Service config, running hook function successfully!")
	} else {
		log.Debugsf("Service config, there hasn't registered any function, skip run hook...")
	}

	return nil
}

// SetConfigs sets the ServiceConfigs into into global variable
func SetConfigs(configs *ServiceConfigs) {
	serviceConfigs = configs
}

// GetConfigs returns the ServiceConfigs from memory
func GetConfigs() *ServiceConfigs {
	return serviceConfigs
}

// GLSCache defines the config of gls cache
type GLSCache struct {
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
	PoolNum  int    `json:"poolNum"`
	Readonly bool   `json:"readonly"`
}

// Equals returns whether the self and other are equals
func (g GLSCache) Equals(o *GLSCache) bool {
	return reflect.DeepEqual(&g, o)
}

func (g GLSCache) String() string {
	return fmt.Sprintf(`GLSCache{Type:%s, Addr:%s, Password:******, PoolNum:%d, Readonly:%++v}`,
		g.Type, g.Addr, g.PoolNum, g.Readonly)
}

// ServiceConfigs defines service-related configuration items,
// and some parameters will automatically set default values if they are not configured in the configuration.
type ServiceConfigs struct {
	ServerAddress               string                                     `json:"serverAddress"`
	Port                        int                                        `json:"port"`
	CallbackPort                int                                        `json:"callbackPort"`
	CommType                    string                                     `json:"commType"`
	Service                     Service                                    `json:"service"`
	ClientSideStatusFSM         bool                                       `json:"clientSideStatusFSM"`
	EnableExecutorLogging       bool                                       `json:"enableExecutorLogging"`
	Log                         Log                                        `json:"log"`
	Version                     int64                                      `json:"version"`
	MaxReadTimeoutMilliseconds  int64                                      `json:"maxReadTimeoutMilliseconds"`
	MaxWriteTimeoutMilliseconds int64                                      `json:"maxWriteTimeoutMilliseconds"`
	Db                          map[string]Db                              `json:"db"`
	Cache                       map[string]Cache                           `json:"cache"`
	Transaction                 Transaction                                `json:"transaction"`
	Heartbeat                   Heartbeat                                  `json:"heartbeat"`
	Alert                       Alert                                      `json:"alert"`
	Apm                         Apm                                        `json:"apm"`
	Deployment                  Deployment                                 `json:"deployment"`
	Addressing                  Addressing                                 `json:"addressing"`
	Downstream                  map[string]Downstream                      `json:"downstream"`
	EventKey                    map[string]string                          `json:"eventKey"`
	GetFn                       func(path string) interface{}              `json:"-"`
	UnmarshalKeyFn              func(key string, rawVal interface{}) error `json:"-"`
	IsExistsFn                  func(path string) bool                     `json:"-"`
	GetIntFn                    func(path string) int                      `json:"-"`
	SetFn                       func(path string, value interface{})       `json:"-"`
	GetInt32Fn                  func(path string) int32                    `json:"-"`
	GetInt64Fn                  func(path string) int64                    `json:"-"`
	GetFloat64Fn                func(path string) float64                  `json:"-"`
	GetStringFn                 func(path string) string                   `json:"-"`
	GetBoolFn                   func(path string) bool                     `json:"-"`
	GetUintFn                   func(path string) uint                     `json:"-"`
	GetUint32Fn                 func(path string) uint32                   `json:"-"`
	GetUint64Fn                 func(path string) uint64                   `json:"-"`
}

// Equals returns whether the self and other are equals
func (s ServiceConfigs) Equals(o *ServiceConfigs) bool {
	return reflect.DeepEqual(&s, o)
}

func (s ServiceConfigs) ToJsonBytes() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ServiceConfigs) RoateFromJsonBytes(source []byte) error {
	if err := json.Unmarshal(source, s); nil != err {
		return err
	}

	return nil
}

func GetMD5(s []byte) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// GetDownstreamServiceConfig returns the downstream service configs Downstream from ServiceConfigs
func (s *ServiceConfigs) GetDownstreamServiceConfig(key string) *Downstream {
	var downStreamServiceConfig Downstream
	if v, ok := s.Downstream[key]; ok {
		downStreamServiceConfig = v
	} else if v, ok := s.Downstream[strings.ToLower(key)]; ok {
		downStreamServiceConfig = v
	}

	return &downStreamServiceConfig
}

func (s *ServiceConfigs) MD5String() string {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return GetMD5(jsonBytes)
}

// Service stores configuration data of [service] section
type Service struct {
	ServiceID                     string                 `json:"serviceID"`
	Org                           string                 `json:"org"`
	Az                            string                 `json:"az"`
	Wks                           string                 `json:"wks"`
	Env                           string                 `json:"env"`
	NodeID                        string                 `json:"nodeID"`
	InstanceID                    string                 `json:"instanceID"`
	Su                            string                 `json:"su"`
	GroupSu                       string                 `json:"groupSu"`
	CommonSu                      string                 `json:"commonSu"`
	ResponseTemplate              string                 `json:"responseTemplate"`
	CustomResponseTemplate        CustomResponseTemplate `json:"customResponseTemplate"`
	CircuitBreakerMonitorPort     int                    `json:"circuitBreakerMonitorPort"`
	ResponseAutoParseKeyMapping   map[string]string      `json:"responseAutoParseKeyMapping"`
	SkipHandleInterceptors        []string               `json:"skipHandleInterceptors"`
	SkipHandleInterceptorsMap     map[string]bool        `json:"skipHandleInterceptorsMap"`
	SkipRemoteCallWrappers        []string               `json:"skipRemoteCallWrappers"`
	SkipRemoteCallWrappersMap     map[string]bool        `json:"skipRemoteCallWrappersMap"`
	SkipRemoteCallInterceptors    []string               `json:"skipRemoteCallInterceptors"`
	SkipRemoteCallInterceptorsMap map[string]bool        `json:"skipRemoteCallInterceptorsMap"`
	ResponseCodeMapping           map[string]string      `json:"responseCodeMapping"`
	DuplicateErrorCodeTo          string                 `json:"duplicateErrorCodeTo"`
}

// CustomResponseTemplate stores the custom response telmpate
type CustomResponseTemplate struct {
	Value string `json:"value"`
}

// Equals returns whether the self and other are equals
func (c CustomResponseTemplate) Equals(o CustomResponseTemplate) bool {
	return reflect.DeepEqual(&c, o)
}

// Equals returns whether the self and other are equals
func (s Service) Equals(o *Service) bool {
	return reflect.DeepEqual(&s, o)
}

// IsMarkedAsSkippedHandleInterceptor checks whether the interceptor is in the skip list or not
func (s Service) IsMarkedAsSkippedHandleInterceptor(interceptorName string) bool {
	if nil != s.SkipHandleInterceptorsMap {
		return s.SkipHandleInterceptorsMap[strings.ToLower(interceptorName)]
	}

	return false
}

// IsMarkedAsSkippedRemoteCallWrapper checks whether the wrapper is in the skip list or not
func (s Service) IsMarkedAsSkippedRemoteCallWrapper(wrapperName string) bool {
	if nil != s.SkipRemoteCallWrappersMap {
		return s.SkipRemoteCallWrappersMap[strings.ToLower(wrapperName)]
	}

	return false
}

// IsMarkedAsSkippedRemoteCallInterceptor checks whether the interceptor is in the skip list or not
func (s Service) IsMarkedAsSkippedRemoteCallInterceptor(interceptorName string) bool {
	if nil != s.SkipRemoteCallInterceptorsMap {
		return s.SkipRemoteCallInterceptorsMap[strings.ToLower(interceptorName)]
	}

	return false
}

// Addressing stores configuration of [addressing] section
type Addressing struct {
	Enable                       bool           `json:"enable"`
	SyncConfigWithServer         bool           `json:"syncConfigWithServer"`
	TopicIDOfServer              string         `json:"topicIDOfServer"`
	TopicVersionOfServer         string         `json:"topicVersionOfServer"`
	TopicSuTitle                 string         `json:"topicSuTitle"`
	DisableGLSLookupOptimization bool           `json:"disableGLSLookupOptimization"`
	RandomElementIDList          []string       `json:"randomElementIDList"`
	RandomElementIDMap           map[string]int `json:"randomElementIDMap"`
	RandomTopicIDList            []string       `json:"randomTopicIDList"`
	RandomTopicIDMap             map[string]int `json:"randomTopicIDMap"`
	Cache                        GLSCache       `json:"cache"`
}

// Equals returns whether the self and other are equals
func (a Addressing) Equals(o *Addressing) bool {
	return reflect.DeepEqual(&a, o)
}

// Log stores configuration data of [log] section
type Log struct {
	LogFile               string   `json:"logFile"`
	LogLevel              string   `json:"logLevel"`
	LogLevelUnixSocket    string   `json:"logLevelUnixSocket"`
	LogFileRootPath       string   `json:"logFileRootPath"`
	MaxSize               int      `json:"maxSize"`
	MaxDays               int      `json:"maxDays"`
	MaxBackups            int      `json:"maxBackups"`
	Console               bool     `json:"console"`
	MaskRules             []string `json:"maskRules"`
	RequestBodyMaskRules  []string `json:"requestBodyMaskRules"`
	ResponseBodyMaskRules []string `json:"responseBodyMaskRules"`

	HeaderMaskRules         []string `json:"headerMaskRules"`
	RequestHeaderMaskRules  []string `json:"requestHeaderMaskRules"`
	ResponseHeaderMaskRules []string `json:"responseHeaderMaskRules"`
}

// Equals returns whether the self and other are equals
func (l Log) Equals(o *Log) bool {
	return reflect.DeepEqual(&l, o)
}

// Apm stores configuration of [apm] section
type Apm struct {
	Enable                              bool   `json:"enable"`
	PrintEmptyTraceIdRecordAtClientSide bool   `json:"printEmptyTraceIdRecordAtClientSide"`
	Version                             string `json:"version"`
	RootPath                            string `json:"rootPath"`
	FileRows                            int    `json:"fileRows"`
}

// Equals returns whether the self and other are equals
func (a Apm) Equals(o *Apm) bool {
	return reflect.DeepEqual(&a, o)
}

// Db stores configuration data of [db] section
type Db struct {
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	Su               string   `json:"su"`
	Topics           []string `json:"topics"`
	Default          bool     `json:"default"`
	Addr             string   `json:"addr"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
	Database         string   `json:"database"`
	Params           string   `json:"params"`
	ServerCACertFile string   `json:"serverCACertFile"`
	ClientCertFile   string   `json:"clientCertFile"`
	ClientPriKeyFile string   `json:"clientPriKeyFile"`
	DBTimeZone       string   `json:"dbTimeZone"`
	Debug            bool     `json:"debug"`
	Pool             struct {
		MaxIdleConns int `json:"maxIdleConns"`
		MaxOpenConns int `json:"maxOpenConns"`
		MaxIdleTime  int `json:"maxIdleTime"`
		MaxLifeValue int `json:"maxLifeValue"`
	} `json:"pool"`
}

// Equals returns whether the self and other are equals
func (d Db) Equals(o *Db) bool {
	return reflect.DeepEqual(&d, o)
}

// EqualsWithoutTopics returns whether the self and other are equals without topics
func (d Db) EqualsWithoutTopics(o *Db) bool {
	return nil != o && d.Type == o.Type &&
		d.Su == o.Su &&
		d.Default == o.Default &&
		d.Addr == o.Addr &&
		d.User == o.User &&
		d.ServerCACertFile == o.ServerCACertFile &&
		d.ClientCertFile == o.ClientCertFile &&
		d.ClientPriKeyFile == o.ClientPriKeyFile &&
		d.Password == o.Password &&
		d.Database == o.Database &&
		d.Params == o.Params &&
		d.Debug == o.Debug &&
		d.DBTimeZone == o.DBTimeZone &&
		d.Pool.MaxIdleConns == o.Pool.MaxIdleConns &&
		d.Pool.MaxOpenConns == o.Pool.MaxOpenConns &&
		d.Pool.MaxIdleTime == o.Pool.MaxIdleTime &&
		d.Pool.MaxLifeValue == o.Pool.MaxLifeValue
}

func (d Db) Clone() Db {
	topics := make([]string, 0)
	for _, topic := range d.Topics {
		topics = append(topics, topic)
	}
	db := Db{
		Name:             d.Name,
		Type:             d.Type,
		Su:               d.Su,
		Topics:           topics,
		Default:          d.Default,
		Addr:             d.Addr,
		User:             d.User,
		Password:         d.Password,
		Database:         d.Database,
		Params:           d.Params,
		ServerCACertFile: d.ServerCACertFile,
		ClientCertFile:   d.ClientCertFile,
		ClientPriKeyFile: d.ClientPriKeyFile,
		Debug:            d.Debug,
		DBTimeZone:       d.DBTimeZone,
		Pool: struct {
			MaxIdleConns int `json:"maxIdleConns"`
			MaxOpenConns int `json:"maxOpenConns"`
			MaxIdleTime  int `json:"maxIdleTime"`
			MaxLifeValue int `json:"maxLifeValue"`
		}{
			MaxIdleConns: d.Pool.MaxIdleConns,
			MaxOpenConns: d.Pool.MaxOpenConns,
			MaxIdleTime:  d.Pool.MaxIdleTime,
			MaxLifeValue: d.Pool.MaxLifeValue,
		},
	}
	return db
}

func (d Db) String() string {
	return fmt.Sprintf(`DB{Name: %s, Type: %s, Su: %s, Topics: %++v, Default: %++v, Addr: %s, User: %s, Password: ******, Database: %s, Params: %s, Debug: %++v, DBTimeZone :%++v, Pool :%++v }`,
		d.Name, d.Type, d.Su, d.Topics, d.Default, d.Addr, d.User, d.Database,
		d.Params, d.Debug, d.DBTimeZone, d.Pool)
}

// Cache stores configuration data of [cache] section
type Cache struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Su       string   `json:"su"`
	Topics   []string `json:"topics"`
	Default  bool     `json:"default"`
	Addr     string   `json:"addr"`
	Password string   `json:"password"`
	Pool     struct {
		PoolSize            int `json:"poolSize"`
		MinIdleConns        int `json:"minIdleConns"`
		MaxConnAgeSeconds   int `json:"maxConnAgeSeconds"`
		PoolTimeoutSeconds  int `json:"poolTimeoutSeconds"`
		IdleTimeoutSeconds  int `json:"idleTimeoutSeconds"`
		DialTimeoutSeconds  int `json:"dialTimeoutSeconds"`
		ReadTimeoutSeconds  int `json:"readTimeoutSeconds"`
		WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
	} `json:"pool"`
}

func (c Cache) String() string {
	return fmt.Sprintf(`DB{Name: %s, Type: %s, Su: %s, Topics: %++v, Default: %++v, Addr: %s, Password: ******, Pool:%++v }`,
		c.Name, c.Type, c.Su, c.Topics, c.Default, c.Addr, c.Pool)
}

// Equals returns whether the self and other are equals
func (c Cache) Equals(o *Cache) bool {
	return reflect.DeepEqual(&c, o)
}

// EqualsWithoutTopics returns whether the self and other are equals without topics
func (c Cache) EqualsWithoutTopics(o *Cache) bool {
	return nil != o && c.Type == o.Type &&
		c.Su == o.Su &&
		c.Default == o.Default &&
		c.Addr == o.Addr &&
		c.Password == o.Password &&
		c.Pool.PoolSize == o.Pool.PoolSize &&
		c.Pool.MinIdleConns == o.Pool.MinIdleConns &&
		c.Pool.MaxConnAgeSeconds == o.Pool.MaxConnAgeSeconds &&
		c.Pool.PoolTimeoutSeconds == o.Pool.PoolTimeoutSeconds &&
		c.Pool.IdleTimeoutSeconds == o.Pool.IdleTimeoutSeconds &&
		c.Pool.DialTimeoutSeconds == o.Pool.DialTimeoutSeconds &&
		c.Pool.ReadTimeoutSeconds == o.Pool.ReadTimeoutSeconds &&
		c.Pool.WriteTimeoutSeconds == o.Pool.WriteTimeoutSeconds
}

func (c Cache) Clone() Cache {
	topics := make([]string, 0)
	for _, topic := range c.Topics {
		topics = append(topics, topic)
	}
	cache := Cache{
		Name:     c.Name,
		Type:     c.Type,
		Su:       c.Su,
		Topics:   topics,
		Default:  c.Default,
		Addr:     c.Addr,
		Password: c.Password,
		Pool: struct {
			PoolSize            int `json:"poolSize"`
			MinIdleConns        int `json:"minIdleConns"`
			MaxConnAgeSeconds   int `json:"maxConnAgeSeconds"`
			PoolTimeoutSeconds  int `json:"poolTimeoutSeconds"`
			IdleTimeoutSeconds  int `json:"idleTimeoutSeconds"`
			DialTimeoutSeconds  int `json:"dialTimeoutSeconds"`
			ReadTimeoutSeconds  int `json:"readTimeoutSeconds"`
			WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
		}{
			PoolSize:            c.Pool.PoolSize,
			MinIdleConns:        c.Pool.MinIdleConns,
			MaxConnAgeSeconds:   c.Pool.MaxConnAgeSeconds,
			PoolTimeoutSeconds:  c.Pool.PoolTimeoutSeconds,
			IdleTimeoutSeconds:  c.Pool.IdleTimeoutSeconds,
			DialTimeoutSeconds:  c.Pool.DialTimeoutSeconds,
			ReadTimeoutSeconds:  c.Pool.ReadTimeoutSeconds,
			WriteTimeoutSeconds: c.Pool.WriteTimeoutSeconds,
		},
	}
	return cache
}

// TransactionServer stores configuration data of [transactionServer] section
type TransactionServer struct {
	Org                    string `json:"org"`
	Wks                    string `json:"wks"`
	Env                    string `json:"env"`
	Su                     string `json:"su"`
	NodeID                 string `json:"nodeID"`
	InstanceID             string `json:"instanceID"`
	TxnBeginEventID        string `json:"txnBeginEventID"`
	TxnJoinEventID         string `json:"txnJoinEventID"`
	TxnEndEventID          string `json:"txnEndEventID"`
	AddressURL             string `json:"addressURL"`
	TxnBeginURLPath        string `json:"txnBeginURLPath"`
	TxnJoinURLPath         string `json:"txnJoinURLPath"`
	TxnEndURLPath          string `json:"txnEndURLPath"`
	MacroServiceAddressURL string `json:"macroServiceAddressURL"`
}

// Equals returns whether the self and other are equals
func (t TransactionServer) Equals(o *TransactionServer) bool {
	return reflect.DeepEqual(&t, o)
}

// Deployment stores configuration data of [deployment] section
type Deployment struct {
	EnableSecure                bool     `json:"enableSecure"`
	TopicTypeOfSKM              string   `json:"topicTypeOfSKM"`
	GetSKMPubKeyTopic           string   `json:"getSKMPubKeyTopic"`
	GetSKMPubKeyTopicVersion    string   `json:"getSKMPubKeyTopicVersion"`
	GetServicesKeysTopic        string   `json:"getServicesKeysTopic"`
	GetServicesKeysTopicVersion string   `json:"getServicesKeysTopicVersion"`
	MaxRetryTimes               int      `json:"maxRetryTimes"`
	RequestTimeoutMilliseconds  int      `json:"requestTimeoutMilliseconds"`
	MaxWaitingMilliseconds      int      `json:"maxWaitingMilliseconds"`
	RetryWaitingMilliseconds    int      `json:"retryWaitingMilliseconds"`
	SystemType                  string   `json:"systemType"`
	CryptoKeyType               string   `json:"cryptoKeyType"`
	Mode                        string   `json:"mode"`
	Padding                     string   `json:"padding"`
	CryptoServiceKey            string   `json:"cryptoServiceKey"`
	CryptoKeyPath               []string `json:"cryptoKeyPath"`
}

// Equals returns whether the self and other are equals
func (d Deployment) Equals(o *Deployment) bool {
	return reflect.DeepEqual(&d, o)
}

// TransactionClient stores configuration data of [transactionClient] section
type TransactionClient struct {
	ConfirmEventID     string `json:"confirmEventID"`
	CancelEventID      string `json:"cancelEventID"`
	ParticipantAddress string `json:"participantAddress"`
	ConfirmAddressURL  string `json:"confirmAddressURL"`
	CancelAddressURL   string `json:"cancelAddressURL"`
}

// Equals returns whether the self and other are equals
func (t TransactionClient) Equals(o *TransactionClient) bool {
	return reflect.DeepEqual(&t, o)
}

// Transaction stores configuration data of [transaction] section
type Transaction struct {
	CommType                      string            `json:"commType"`
	IsPropagator                  bool              `json:"isPropagator"`
	TryFailedIgnoreCallbackCancel bool              `json:"tryFailedIgnoreCallbackCancel"`
	SaveHeaders                   bool              `json:"saveHeaders"`
	TimeoutMilliseconds           int               `json:"timeoutMilliseconds"`
	MaxRetryTimes                 int               `json:"maxRetryTimes"`
	MaxServiceConsumeMilliseconds int               `json:"maxServiceConsumeMilliseconds"`
	PropagatorServices            []string          `json:"propagatorServices"`
	PropagatorServicesMap         map[string]bool   `json:"propagatorServicesMap"`
	IsMacroService                bool              `json:"isMacroService"`
	TransactionServer             TransactionServer `json:"transactionServer"`
	TransactionClient             TransactionClient `json:"transactionClient"`
}

// Equals returns whether the self and other are equals
func (t Transaction) Equals(o *Transaction) bool {
	return reflect.DeepEqual(&t, o)
}

// Heartbeat stores configuration data of [heartbeat] section
type Heartbeat struct {
	TopicName       string `json:"topicName"`
	IntervalSeconds int    `json:"intervalSeconds"`
}

// Equals returns whether the self and other are equals
func (h Heartbeat) Equals(o *Heartbeat) bool {
	return reflect.DeepEqual(&h, o)
}

// Alert stores configuration data of [alert] section
type Alert struct {
	TopicName string `json:"topicName"`
}

// Equals returns whether the self and other are equals
func (a Alert) Equals(o *Alert) bool {
	return reflect.DeepEqual(&a, o)
}

// Downstream stores configuration data of [downstream] section
type Downstream struct {
	EventType                        string               `json:"eventType"`
	EventID                          string               `json:"eventID"`
	Version                          string               `json:"version"`
	TimeoutMilliseconds              int                  `json:"timeoutMilliseconds"`
	RetryWaitingMilliseconds         int                  `json:"retryWaitingMilliseconds"`
	MaxWaitingTimeMilliseconds       int                  `json:"maxWaitingTimeMilliseconds"`
	MaxRetryTimes                    int                  `json:"maxRetryTimes"`
	DeleteTransactionPropagationInfo bool                 `json:"deleteTransactionPropagationInfo"`
	ProtoType                        string               `json:"protoType"`
	HTTPAddress                      string               `json:"httpAddress"`
	HTTPMethod                       string               `json:"httpMethod"`
	HTTPContextType                  string               `json:"httpContextType"`
	ResponseAutoParseKeyMapping      map[string]string    `json:"responseAutoParseKeyMapping"`
	PassThroughHeaderKey             PassThroughHeaderKey `json:"ResponseAutoParseKeyMapping"`
	CircuitBreaker                   CircuitBreaker       `json:"circuitBreaker"`
	CustomConfigurations             CustomConfigurations `json:"customConfigurations"`
	EnableLogging                    bool                 `json:"enableLogging"`
	//Masker                           Masker               `json:"masker"`
}

// Masker stores configuration of [downstream.XXXXX.masker] section
type Masker struct {
	MaskRules             []string `json:"maskRules"`
	RequestBodyMaskRules  []string `json:"requestBodyMaskRules"`
	ResponseBodyMaskRules []string `json:"responseBodyMaskRules"`

	HeaderMaskRules         []string `json:"headerMaskRules"`
	RequestHeaderMaskRules  []string `json:"requestHeaderMaskRules"`
	ResponseHeaderMaskRules []string `json:"responseHeaderMaskRules"`
}

// CustomConfigurations stores configuration of [CustomConfigurations] section
type CustomConfigurations struct {
	TimeoutMilliseconds              int    `json:"timeoutMilliseconds"`
	RetryWaitingMilliseconds         int    `json:"retryWaitingMilliseconds"`
	MaxWaitingTimeMilliseconds       int    `json:"maxWaitingTimeMilliseconds"`
	MaxRetryTimes                    int    `json:"maxRetryTimes"`
	DeleteTransactionPropagationInfo bool   `json:"deleteTransactionPropagationInfo"`
	ProtoType                        string `json:"protoType"`
	HTTPAddress                      string `json:"httpAddress"`
	HTTPMethod                       string `json:"httpMethod"`
	HTTPContextType                  string `json:"httpContextType"`
}

// PassThroughHeaderKey stores configuration of [PassThroughHeaderKey] section
type PassThroughHeaderKey struct {
	List []string `json:"list"`
}

// CircuitBreaker stores configuration of [circuitBreaker] section
type CircuitBreaker struct {
	Enable                  bool `json:"enable"`
	MaxConcurrentRequests   int  `json:"maxConcurrentRequests"`
	SleepWindowMilliseconds int  `json:"sleepWindowMilliseconds"`
	RequestVolumeThreshold  int  `json:"requestVolumeThreshold"`
	ErrorPercentThreshold   int  `json:"errorPercentThreshold"`
}

// Equals returns whether the self and other are equals
func (d Downstream) Equals(o *Downstream) bool {
	return reflect.DeepEqual(&d, o)
}

// GetDownstreamIgnoreCaseKey finds the value of downstream by key and ignore case the key.
func GetDownstreamIgnoreCaseKey(m map[string]Downstream, key string) *Downstream {
	var downStreamServiceConfig Downstream
	if v, ok := m[key]; ok {
		downStreamServiceConfig = v
	} else if v, ok := m[strings.ToLower(key)]; ok {
		downStreamServiceConfig = v
	}

	return &downStreamServiceConfig
}

// GenLogConfig generates a log config(log.Config) from ServiceConfigs
func (s *ServiceConfigs) GenLogConfig() log.Config {
	return log.Config{
		Rotate: log.RotateConfig{
			FilePath:   s.Log.LogFileRootPath,
			Filename:   s.Log.LogFile,
			MaxSize:    s.Log.MaxSize,
			MaxBackups: s.Log.MaxBackups,
			MaxAge:     s.Log.MaxDays,
		},
		Level:   s.Log.LogLevel,
		Console: s.Log.Console,
	}
}

// GenCallbackOptions generates some callback.Option from ServiceConfigs
func (s *ServiceConfigs) GenCallbackOptions() []callback.Option {
	return []callback.Option{
		func(options *callback.Options) {
			options.Port = s.Port
			options.CallbackPort = s.CallbackPort
			options.ServerAddress = s.ServerAddress
			options.CommType = s.CommType
			options.EnableClientSideStatusFSM = s.ClientSideStatusFSM
			extConfigMap := make(map[string]interface{}, 0)
			extConfigMap[constant.ExtConfigTransaction] = &s.Transaction
			extConfigMap[constant.ExtConfigDownstreamService] = s.Downstream
			extConfigMap[constant.ExtEventKeyMap] = s.EventKey
			extConfigMap[constant.ExtConfigService] = &s.Service
			extConfigMap[constant.ExtEnableExecutorLogging] = s.EnableExecutorLogging
			options.ExtConfigs = extConfigMap
		},
	}
}

// Get gets the interface{} value from ServiceConfigs by the key
func (s *ServiceConfigs) Get(key string) interface{} {
	return s.GetFn(key)
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func (s *ServiceConfigs) UnmarshalKey(key string, rawVal interface{}) error {
	return s.UnmarshalKeyFn(key, rawVal)
}

// GetBool gets the boolean value from ServiceConfigs by the key
func (s *ServiceConfigs) GetBool(key string) bool {
	return s.GetBoolFn(key)
}

// DefaultBool gets the boolean value from ServiceConfigs by the key,
// returns defaultVal(bool) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultBool(key string, defaultVal bool) bool {
	if s.IsExistsFn(key) {
		return s.GetBool(key)
	}

	return defaultVal
}

// GetString gets the string value from ServiceConfigs by the key
func (s *ServiceConfigs) GetString(key string) string {
	return s.GetStringFn(key)
}

// DefaultString gets the string value from ServiceConfigs by the key
// returns defaultVal(string) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultString(key string, defaultVal string) string {
	if s.IsExistsFn(key) {
		return s.GetString(key)
	}

	return defaultVal
}

// GetInt gets the int value from ServiceConfigs by the key
func (s *ServiceConfigs) GetInt(key string) int {
	return s.GetIntFn(key)
}

// DefaultInt gets the int value from ServiceConfigs by the key
// returns defaultVal(int) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultInt(key string, defaultVal int) int {
	if s.IsExistsFn(key) {
		return s.GetInt(key)
	}

	return defaultVal
}

// GetInt32 gets the int32 value from ServiceConfigs by the key
func (s *ServiceConfigs) GetInt32(key string) int32 {
	return s.GetInt32Fn(key)
}

// DefaultInt32 gets the int32 value from ServiceConfigs by the key
// returns defaultVal(int32) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultInt32(key string, defaultVal int32) int32 {
	if s.IsExistsFn(key) {
		return s.GetInt32(key)
	}

	return defaultVal
}

// GetInt64 gets the int64 value from ServiceConfigs by the key
func (s *ServiceConfigs) GetInt64(key string) int64 {
	return s.GetInt64Fn(key)
}

// DefaultInt64 gets the int64 value from ServiceConfigs by the key
// returns defaultVal(int64) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultInt64(key string, defaultVal int64) int64 {
	if s.IsExistsFn(key) {
		return s.GetInt64(key)
	}

	return defaultVal
}

// GetFloat64 gets the float64 value from ServiceConfigs by the key
func (s *ServiceConfigs) GetFloat64(key string) float64 {
	return s.GetFloat64Fn(key)
}

// DefaultFloat64 gets the float64 value from ServiceConfigs by the key
// returns defaultVal(float64) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultFloat64(key string, defaultVal float64) float64 {
	if s.IsExistsFn(key) {
		return s.GetFloat64(key)
	}

	return defaultVal
}

// GetUint gets the uint value from ServiceConfigs by the key
func (s *ServiceConfigs) GetUint(key string) uint {
	return s.GetUintFn(key)
}

// DefaultUnit gets the uint value from ServiceConfigs by the key
// returns defaultVal(uint) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultUnit(key string, defaultVal uint) uint {
	if s.IsExistsFn(key) {
		return s.GetUint(key)
	}

	return defaultVal
}

// GetUint32 gets the uint32 value from ServiceConfigs by the key
func (s *ServiceConfigs) GetUint32(key string) uint32 {
	return s.GetUint32Fn(key)
}

// DefaultUnit32 gets the uint32 value from ServiceConfigs by the key
// returns defaultVal(uint32) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultUnit32(key string, defaultVal uint32) uint32 {
	if s.IsExistsFn(key) {
		return s.GetUint32(key)
	}

	return defaultVal
}

// GetUint64 gets the uint64 value from ServiceConfigs by the key
func (s *ServiceConfigs) GetUint64(key string) uint64 {
	return s.GetUint64Fn(key)
}

// DefaultUnit64 gets the uint64 value from ServiceConfigs by the key
// returns defaultVal(uint64) if the key path doesn't exist in the ServiceConfigs
func (s *ServiceConfigs) DefaultUnit64(key string, defaultVal uint64) uint64 {
	if s.IsExistsFn(key) {
		return s.GetUint64(key)
	}

	return defaultVal
}

// GetEventID finds the event ID in the `event key`-`event ID` map by `event key`
func GetEventID(m map[string]string, key string) string {
	eventID := ""
	if v, ok := m[key]; ok {
		eventID = v
	} else if v, ok := m[strings.ToLower(key)]; ok {
		eventID = v
	}

	return eventID
}
