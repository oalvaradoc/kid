package v2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/logging"
	"git.multiverse.io/eventkit/kit/common/crypto/aes"
	"git.multiverse.io/eventkit/kit/common/crypto/rsa"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/client_receive"
	"git.multiverse.io/eventkit/kit/log"
	masker2 "git.multiverse.io/eventkit/kit/masker/json"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Loader is a config loader using viper to load the config file
type Loader struct{}

const defaultCustomEventIDRelationFilePath = "conf/custom_event_id_relation.toml"

var currentConfigMD5 string

func IsFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
}

func ReadFileMd5(sfile string) (string, error) {
	ssconfig, err := ioutil.ReadFile(sfile)
	if err != nil {
		return "", err
	}
	return config.GetMD5(ssconfig), nil
}

// LoadConfig performs the config file load logics
func (b *Loader) LoadConfig(filePath string) (*config.ServiceConfigs, error) {
	isExistsCustomEventIDRelationFile := false
	var customEventIDRelationSetting map[string]interface{}
	if IsFileExists(defaultCustomEventIDRelationFilePath) {
		isExistsCustomEventIDRelationFile = true
		viper.SetConfigFile(defaultCustomEventIDRelationFilePath)
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}

		customEventIDRelationSetting = viper.AllSettings()
		viper.Reset()
	}


	viper.SetConfigFile(filePath)
	viper.SetDefault("CommType", "mesh")
	viper.SetDefault("ServerAddress", "http://127.0.0.1:18080")
	viper.SetDefault("CallbackPort", 18082)
	viper.SetDefault("Port", 6060)
	viper.SetDefault("maxReadTimeoutMilliseconds", 80*1000)
	viper.SetDefault("maxWriteTimeoutMilliseconds", 80*1000)
	viper.SetDefault("service.circuitBreakerMonitorPort", constant.DefaultCircuitBreakerMonitorPort)
	viper.SetDefault("log.console", true)
	viper.SetDefault("log.logLevel", "info")
	viper.SetDefault("log.filePath", "/data/logs")
	viper.SetDefault("log.filename", "default.log")
	viper.SetDefault("log.maxSize", 200)
	viper.SetDefault("log.maxBackups", 7)
	viper.SetDefault("log.maxAge", 7)

	viper.SetDefault("transaction.tryFailedIgnoreCallbackCancel", true)
	viper.SetDefault("transaction.transactionServer.txnBeginURLPath", constant.TxnBeginURLPath)
	viper.SetDefault("transaction.transactionServer.txnJoinURLPath", constant.TxnJoinURLPath)
	viper.SetDefault("transaction.transactionServer.txnEndURLPath", constant.TxnEndURLPath)
	viper.SetDefault("transaction.transactionServer.macroServiceAddressURL", constant.TxnMacroServiceAddressURL)

	viper.SetDefault("addressing.topicSuTitle", "TOP.GLSTOPIC")
	viper.SetDefault("addressing.topicIDOfServer", "GlsAppConfig")

	viper.SetDefault("apm.rootPath", "/data/logs/")
	viper.SetDefault("apm.version", "v2")
	viper.SetDefault("apm.fileRows", 1000000)

	viper.SetDefault("deployment.getSKMPubKeyTopic", "GetSkmSSLPem")
	viper.SetDefault("deployment.getServicesKeysTopic", "GetKeysWithCrypto")
	viper.SetDefault("deployment.maxRetryTimes", 3)
	viper.SetDefault("deployment.requestTimeoutMilliseconds", 30000)
	viper.SetDefault("deployment.maxWaitingMilliseconds", 120000)
	viper.SetDefault("deployment.retryWaitingMilliseconds", 10000)
	viper.SetDefault("deployment.cryptoKeyType", "RSA2048")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if isExistsCustomEventIDRelationFile {
		allSettings := viper.AllSettings()
		viper.Reset()
		viper.SetConfigFile(filePath)

		delete(allSettings, "downstream")
		delete(allSettings, "eventKey")
		viper.MergeConfigMap(allSettings)
		viper.MergeConfigMap(customEventIDRelationSetting)
	}


	cfg := &config.ServiceConfigs{}
	viper.Unmarshal(cfg)

	completingOtherConfigurations(cfg)
	callback.SetMaxReadAndWriteTimeoutMilliseconds(cfg.MaxReadTimeoutMilliseconds, cfg.MaxWriteTimeoutMilliseconds)
	config.SetConfigs(cfg)
	currentConfigMD5, err = ReadFileMd5(filePath)
	if nil != err {
		return nil, err
	}
	callback.SetConfigVersion(cfg.Version)
	// register on config change event function to watch service config.
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		currentConfig := config.GetConfigs()
		fileMD5, err := ReadFileMd5(filePath)
		if nil != err {
			log.Errorsf("Failed to read file MD5[event=%++v], error:%++v", e, err)
		}
		if fileMD5 == currentConfigMD5 {
			log.Debugsf("config file no change, skip roate service config")
			return
		}
		currentConfigMD5 = fileMD5
		log.Infosf("file config change, start change service config...")

		err = viper.ReadInConfig()
		if err != nil {
			log.Errorsf("Failed to read service config, error:%++v", err)
			return
		}
		if isExistsCustomEventIDRelationFile {
			allSettings := viper.AllSettings()
			viper.Reset()
			viper.SetConfigFile(filePath)

			delete(allSettings, "downstream")
			delete(allSettings, "eventKey")
			viper.MergeConfigMap(allSettings)
			viper.MergeConfigMap(customEventIDRelationSetting)
		}
		newCfg := &config.ServiceConfigs{}
		viper.Unmarshal(newCfg)

		completingOtherConfigurations(newCfg)
		callback.SetMaxReadAndWriteTimeoutMilliseconds(newCfg.MaxReadTimeoutMilliseconds, newCfg.MaxWriteTimeoutMilliseconds)
		config.SetConfigs(newCfg)
		// must after the `completingOtherConfigurations`
		config.RunHookHandleFuncIfNecessary(currentConfig, newCfg)
		log.Debugsf("new service config:%++v", newCfg)

		callback.SetConfigVersion(newCfg.Version)
		log.Infosf("Successfully to updated service config on config file changed!")
	})

	return cfg, nil
}

func completingOtherConfigurations(cfg *config.ServiceConfigs) {
	cfg.GetFn = Get
	cfg.IsExistsFn = IsExistsFn
	cfg.GetIntFn = GetInt
	cfg.GetInt32Fn = GetInt32
	cfg.SetFn = Set
	cfg.GetInt64Fn = GetInt64
	cfg.GetFloat64Fn = GetFloat64
	cfg.GetStringFn = GetString
	cfg.GetBoolFn = GetBool
	cfg.GetUintFn = GetUint
	cfg.GetUint32Fn = GetUint32
	cfg.GetUint64Fn = GetUint64
	cfg.UnmarshalKeyFn = UnmarshalKey

	// decrypt the config if necessary
	decryptTheConfigIfNecessary(cfg)

	if strings.HasPrefix(cfg.Service.NodeID, "$") {
		nodeID := os.Getenv("NODE_ID")
		if "" == nodeID {
			panic("cannot found env key `NODE_ID`")
		}
		log.Infosf("Replace cfg.Service.NodeID with:[%s]", nodeID)
		cfg.Service.NodeID = nodeID
	}

	if strings.HasPrefix(cfg.Transaction.TransactionServer.NodeID, "$") {
		nodeID := os.Getenv("NODE_ID")
		if "" == nodeID {
			panic("cannot found env key `NODE_ID`")
		}
		log.Infosf("Replace cfg.Transaction.TransactionServer.NodeID with:[%s]", nodeID)
		cfg.Transaction.TransactionServer.NodeID = nodeID
	}

	if strings.HasPrefix(cfg.Transaction.TransactionServer.InstanceID, "$") {
		instanceID := os.Getenv("INSTANCE_ID")
		if "" == instanceID {
			panic("cannot found env key `INSTANCE_ID`")
		}
		log.Infosf("Replace cfg.Transaction.TransactionServer.InstanceID with:[%s]", instanceID)
		cfg.Transaction.TransactionServer.InstanceID = instanceID
	}

	cfg.Service.InstanceID = strings.ReplaceAll(strings.ReplaceAll(os.Getenv("INSTANCE_ID"), "-", ""), "_", "")

	participantAddress := fmt.Sprintf("http://%s:%d", CurrentHost(), cfg.Port)

	cfg.Transaction.TransactionClient.ParticipantAddress = participantAddress
	cfg.Transaction.PropagatorServicesMap = make(map[string]bool)
	if len(cfg.Transaction.PropagatorServices) > 0 {
		for _, serviceName := range cfg.Transaction.PropagatorServices {
			cfg.Transaction.PropagatorServicesMap[serviceName] = true
		}
	}

	cfg.Service.SkipHandleInterceptorsMap = make(map[string]bool)
	if len(cfg.Service.SkipHandleInterceptors) > 0 {
		for _, interceptorName := range cfg.Service.SkipHandleInterceptors {
			cfg.Service.SkipHandleInterceptorsMap[strings.ToLower(interceptorName)] = true
		}
	}

	cfg.Service.SkipRemoteCallInterceptorsMap = make(map[string]bool)
	if len(cfg.Service.SkipRemoteCallInterceptors) > 0 {
		for _, interceptorName := range cfg.Service.SkipRemoteCallInterceptors {
			cfg.Service.SkipRemoteCallInterceptorsMap[strings.ToLower(interceptorName)] = true
		}
	}

	cfg.Service.SkipRemoteCallWrappersMap = make(map[string]bool)
	if len(cfg.Service.SkipRemoteCallWrappers) > 0 {
		for _, wrapperName := range cfg.Service.SkipRemoteCallWrappers {
			cfg.Service.SkipRemoteCallWrappersMap[strings.ToLower(wrapperName)] = true
		}
	}

	cfg.Addressing.RandomElementIDMap = make(map[string]int)
	if len(cfg.Addressing.RandomElementIDList) > 0 {
		for _, wrapperName := range cfg.Addressing.RandomElementIDList {
			keyArray := strings.Split(wrapperName, ":")
			if len(keyArray) > 1 {
				if v, err := strconv.Atoi(keyArray[1]); nil != err {
					panic(errors.Errorf(constant.SystemInternalError, "Failed to parse string to int for element_id_list, key = %s, error = %++v", wrapperName, err))
				} else {
					if v <= 0 {
						panic(errors.Errorf(constant.SystemInternalError, "Failed to parse string to int for element_id_list, sharedNumberString cannot be zero, key = %s", wrapperName))
					}
					cfg.Addressing.RandomElementIDMap[strings.ToLower(keyArray[0])] = v
				}
			} else {
				cfg.Addressing.RandomElementIDMap[strings.ToLower(wrapperName)] = 10
			}
		}
	}

	cfg.Addressing.RandomTopicIDMap = make(map[string]int)
	if len(cfg.Addressing.RandomTopicIDList) > 0 {
		for _, wrapperName := range cfg.Addressing.RandomTopicIDList {
			keyArray := strings.Split(wrapperName, ":")
			if len(keyArray) > 1 {
				if v, err := strconv.Atoi(keyArray[1]); nil != err {
					panic(errors.Errorf(constant.SystemInternalError, "Failed to parse string to int for topic_id_list, key = %s, error = %++v", wrapperName, err))
				} else {
					if v <= 0 {
						panic(errors.Errorf(constant.SystemInternalError, "Failed to parse string to int for topic_id_list, sharedNumberString cannot be zero, key = %s", wrapperName))
					}
					cfg.Addressing.RandomTopicIDMap[strings.ToLower(keyArray[0])] = v
				}
			} else {
				cfg.Addressing.RandomTopicIDMap[strings.ToLower(wrapperName)] = 10
			}
		}
	}

	errors.SetErrorCodeMapping(cfg.Service.ResponseCodeMapping)
}

var (
	cryptoPriKey                 []byte = nil
	cryptoPubKey                 []byte = nil
	cryptoAesKey                 []byte = nil
	cryptoKeyLoaderOnce          sync.Once
	defaultCallWrapperOption     = client.DefaultWrapperCall(&logging.Wrapper{})
	defaultCallInterceptorOption = client.DefaultCallInterceptors([]interceptor.Interceptor{&client_receive.Interceptor{}})
)

const (
	ConfigDecryptionMaskerName = "configDecryptionMasker"
	SkmPublicPemAlgorithm      = "rsa"
)

type ConfigDecryptionMasker struct{
	CryptoKeyType string
	Mode string
	Padding string
}

func (p *ConfigDecryptionMasker) Do(keyPath, in string, _ ...string) (string, error) {
	var result []byte
	if strings.HasPrefix(strings.ToLower(p.CryptoKeyType), "aes") {
		// decrypt the source parameter string using AES
		if nil == cryptoAesKey {
			return "", errors.Errorf(constant.SystemInternalError, "Can not get the AES key for decryption masker, key path:[%s]", keyPath)
		}

		aes := aes.NewAES(p.Padding, p.Mode, cryptoAesKey)

		encryptedBytes, err := base64.StdEncoding.DecodeString(in)
		if nil != err {
			return "", errors.Errorf(constant.SystemInternalError, "Failed to decode input string, key path:[%s], error:%++v", keyPath, err)
		}
		result, err = aes.Decrypt(encryptedBytes)
		if nil != err {
			return "", errors.Errorf(constant.SystemInternalError, "Failed to decrypt input string, key path:[%s], error:%++v", keyPath, err)
		}
	} else {
		// decrypt the source parameter string using RSA
		if nil == cryptoPriKey {
			return "", errors.Errorf(constant.SystemInternalError, "Can not get the RSA private key for decryption masker, key path:[%s]", keyPath)
		}
		rsa := rsa.NewRSA(cryptoPriKey, cryptoPubKey)
		encryptedBytes, err := base64.StdEncoding.DecodeString(in)
		if nil != err {
			return "", errors.Errorf(constant.SystemInternalError, "Failed to decode input string, key path:[%s], error:%++v", keyPath, err)
		}
		result, err = rsa.Decrypt(encryptedBytes)
		if nil != err {
			return "", errors.Errorf(constant.SystemInternalError, "Failed to decrypt input string, key path:[%s], error:%++v", keyPath, err)
		}
	}
	maskedValue := string(result)

	viper.Set(keyPath, maskedValue)
	return maskedValue, nil
}

func GetSKMPubKey(client client.Client, cfg *config.ServiceConfigs) (pubKey []byte, err error) {
	var request = SkmPublicPemRequest{
		Algorithm: SkmPublicPemAlgorithm,
	}

	topicTypeOption := mesh.WithTopicTypeBusiness()
	if len(cfg.Deployment.TopicTypeOfSKM) > 0 {
		topicTypeOption = mesh.WithTopicType(cfg.Deployment.TopicTypeOfSKM)
	}

	meshRequest := mesh.NewMeshRequest(request)
	meshRequest.WithOptions(
		topicTypeOption,                              			   // setting topic type
		mesh.WithORG(cfg.Service.Org),                             // org id
		mesh.WithWorkspace(cfg.Service.Wks),                       // workspace
		mesh.WithEnvironment(cfg.Service.Env),                     // environment
		mesh.WithSU(cfg.Service.Su),                               // su
		mesh.WithVersion(cfg.Deployment.GetSKMPubKeyTopicVersion), //version
		mesh.WithEventID(cfg.Deployment.GetSKMPubKeyTopic),        // destination event id
		mesh.WithMaxRetryTimes(cfg.Deployment.MaxRetryTimes),      // retry times
		mesh.WithRetryWaitingMilliseconds(
			time.Duration(cfg.Deployment.RetryWaitingMilliseconds)*time.Millisecond), // retry waiting time
		mesh.WithMaxWaitingTime(
			time.Duration(cfg.Deployment.MaxWaitingMilliseconds)*time.Millisecond), // max waiting time
		mesh.WithTimeout(
			time.Duration(cfg.Deployment.RequestTimeoutMilliseconds)*time.Millisecond), // request timeout milliseconds
	)

	response := &SkmPublicPemResponseData{}
	// sync call
	handlerContext := contexts.BuildHandlerContexts(contexts.ResponseAutoParseKeyMapping(map[string]string{
		"type":            "json",
		"errorCodeKey":    "errorCode",
		"errorMsgKey":     "errorMsg",
		"responseDataKey": "response",
	}))
	ctx := contexts.BuildContextFromParentWithHandlerContexts(context.Background(), handlerContext)
	_, err = client.SyncCall(ctx, meshRequest, response)
	if nil != err {
		return nil, err
	}

	return []byte(response.PublicPem), err
}

func EncryptRequest(request []byte, skmPubKey []byte) (aesKey []byte, krc *SkmRequestCrypto, err error) {
	// generate aes key (32 bytes)
	aesKey = make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, nil, err
	}
	// encrypt request
	aesIns := aes.NewAES("PKCS5", "GCM", aesKey)
	encryptedRequest, err := aesIns.Encrypt(request)

	if nil != err {
		return nil, nil, err
	}

	skmRequestCrypto := &SkmRequestCrypto{}
	skmRequestCrypto.Request = base64.StdEncoding.EncodeToString(encryptedRequest)

	rsaIns := rsa.NewRSA(nil, skmPubKey)
	encryptedKey, err := rsaIns.Encrypt(aesKey)

	if nil != err {
		return nil, nil, err
	}

	skmRequestCrypto.Key = base64.StdEncoding.EncodeToString(encryptedKey)

	return aesKey, skmRequestCrypto, nil
}

func DecryptResponse(response []byte, aesKey []byte) (decryptedResponse []byte, err error) {
	// decrypt request
	aesIns := aes.NewAES("PKCS5", "GCM", aesKey)

	decryptedResponse, err = aesIns.Decrypt(response)

	return decryptedResponse, err
}

var cryptoPubKeyTypeSuffix = "-PUB"
var cryptoPriKeyTypeSuffix = "-PRI"

func loadCryptoKeyFromSKM(cfg *config.ServiceConfigs) {
	cryptoKeyLoaderOnce.Do(func() {
		callback.SetSedServerAddr(cfg.ServerAddress)
		client := mesh.NewMeshClient(defaultCallWrapperOption, defaultCallInterceptorOption)
		// 1. get skm communicate public key
		skmPubKey, gerr := GetSKMPubKey(client, cfg)
		if nil != gerr {
			log.Errorsf("Failed to get skm public key, error:%++v", gerr)
			return
		}

		// 2. using skm communtion public key and crypto service key to get RSA private key
		var body []byte
		var err error

		services := make([]SkmServiceRequest, 0)
		effectTime := time.Now()
		systemType := "APP"
		if len(cfg.Deployment.SystemType) > 0 {
			systemType = cfg.Deployment.SystemType
		}
		services = append(services, SkmServiceRequest{
			KeyType:    KeyType(cfg.Deployment.CryptoKeyType),
			ServiceId:  cfg.Deployment.CryptoServiceKey,
			SystemType: systemType,
			Operator: cfg.Service.ServiceID + "-" +
				cfg.Service.InstanceID + "-" + cfg.Service.NodeID,
			RequestTime: Date(effectTime),
		})

		request := &SkmKeyValuesRequest{
			Services: services,
		}
		body, err = json.Marshal(request)

		if err != nil {
			log.Errorsf("Failed to marshal request into json string, error:%++v", err)
			return
		}
		log.Debugsf("loadCryptoKeyFromSKM service request:[%s]", string(body))

		aesKey, skmRequestCrypto, err := EncryptRequest(body, skmPubKey)
		if err != nil {
			log.Errorsf("Failed to encrypt SKM request, error:%++v", err)
			return
		}

		topicTypeOption := mesh.WithTopicTypeBusiness()
		if len(cfg.Deployment.TopicTypeOfSKM) > 0 {
			topicTypeOption = mesh.WithTopicType(cfg.Deployment.TopicTypeOfSKM)
		}

		meshRequest := mesh.NewMeshRequest(skmRequestCrypto)
		meshRequest.WithOptions(
			topicTypeOption,                              			   	  // setting topic type
			mesh.WithORG(cfg.Service.Org),                                // org id
			mesh.WithWorkspace(cfg.Service.Wks),                          // workspace
			mesh.WithEnvironment(cfg.Service.Env),                        // environment
			mesh.WithSU(cfg.Service.Su),                                  // su
			mesh.WithVersion(cfg.Deployment.GetServicesKeysTopicVersion), // version
			mesh.WithEventID(cfg.Deployment.GetServicesKeysTopic),        // destination event id
			mesh.WithMaxRetryTimes(cfg.Deployment.MaxRetryTimes),         // retry times
			mesh.WithRetryWaitingMilliseconds(
				time.Duration(cfg.Deployment.RetryWaitingMilliseconds)*time.Millisecond), // retry waiting time
			mesh.WithMaxWaitingTime(
				time.Duration(cfg.Deployment.MaxWaitingMilliseconds)*time.Millisecond), // max waiting time
			mesh.WithTimeout(
				time.Duration(cfg.Deployment.RequestTimeoutMilliseconds)*time.Millisecond), // request timeout milliseconds
		)

		response := &SkmResponseCrypto{}
		// sync call
		handlerContext := contexts.BuildHandlerContexts(contexts.ResponseAutoParseKeyMapping(map[string]string{
			"type":            "json",
			"errorCodeKey":    "errorCode",
			"errorMsgKey":     "errorMsg",
			"responseDataKey": "response",
		}))

		ctx := contexts.BuildContextFromParentWithHandlerContexts(context.Background(), handlerContext)
		_, err = client.SyncCall(ctx, meshRequest, response)
		if nil != err {
			log.Errorsf("Failed to sync call SKM to get RSA primary key, error:%++v", err)
			return
		}

		skmKeysResponse, ugerr := unpackResponseMsg(aesKey, response)
		if nil != ugerr {
			log.Errorsf("Failed to unpack response message into SKM keys response struct. error:%++v", ugerr)
			return
		}

		if len(skmKeysResponse.Result) > 0 {
			firstKeyResult := skmKeysResponse.Result[0]
			for _, e := range firstKeyResult.Result {
				if e.KeyType == KeyType(cfg.Deployment.CryptoKeyType+cryptoPubKeyTypeSuffix) {
					// public key
					cryptoPubKey = []byte(e.KeyValue)
				}

				if e.KeyType == KeyType(cfg.Deployment.CryptoKeyType+cryptoPriKeyTypeSuffix) {
					// private key
					cryptoPriKey = []byte(e.KeyValue)
				}

				if e.KeyType == KeyType(cfg.Deployment.CryptoKeyType) {
					// aes key
					cryptoAesKey = []byte(e.KeyValue)
				}
			}
		}

	})
}

func unpackResponseMsg(aesKey []byte, skmResponseCrypto *SkmResponseCrypto) (*SkmKeysResponse, *errors.Error) {
	ciphertextResponse, err := base64.StdEncoding.DecodeString(skmResponseCrypto.Response)
	if nil != err {
		log.Errorsf("unpackResponseMsg:unpack response message failed, error:[%s], skmResponseCrypto.Response:[%++v]",
			errors.ErrorToString(err), skmResponseCrypto.Response)

		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	aesIns := aes.NewAES("PKCS5", "GCM", aesKey)

	// decrypt reply body
	decryptedResponse, err := aesIns.Decrypt(ciphertextResponse)

	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	result := SkmKeysResponse{}
	err = json.Unmarshal(decryptedResponse, &result)
	if err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return &result, nil
}

func decryptTheConfigIfNecessary(cfg *config.ServiceConfigs) {
	if !cfg.Deployment.EnableSecure {
		log.Infosf("Secure deployment is disable, skip decryption...")
		return
	}

	if len(cfg.Deployment.CryptoKeyPath) == 0 {
		log.Infosf("Crypto key path is empty, skip decryption...")
		return
	}
	masker2.RegisterMasker(ConfigDecryptionMaskerName, &ConfigDecryptionMasker{
		CryptoKeyType: cfg.Deployment.CryptoKeyType,
		Mode:          cfg.Deployment.Mode,
		Padding:       cfg.Deployment.Padding,
	})

	loadCryptoKeyFromSKM(cfg)

	originalBytes, gerr := json.Marshal(viper.AllSettings())
	if nil != gerr {
		panic(gerr)
	}
	maskerRuleKeys := make([]string, 0)
	for _, cryptoKeyPath := range cfg.Deployment.CryptoKeyPath {
		maskerRuleKeys = append(maskerRuleKeys, strings.ToLower(cryptoKeyPath) +"|" + ConfigDecryptionMaskerName)
	}

	finalConfigJsonBytes, isExistError, gerr := masker2.JsonBodyMask(originalBytes, maskerRuleKeys, true)
	if isExistError || nil != gerr {
		if nil != gerr {
			log.Errorsf("The error of json body mask,error:%++v", gerr)
			panic(gerr)
		}

		panic("Failed to body mask, please check!")
	}

	gerr = cfg.RoateFromJsonBytes(finalConfigJsonBytes)

	if nil != gerr {
		panic(gerr)
	}
}

// CurrentHost gets the current host IP, returns `localhost` if failed to get host IP.
func CurrentHost() (host string) {
	host = "localhost"
	netInterfaces, e := net.Interfaces()
	if e != nil {
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						host = ipnet.IP.String()
						return
					}
				}
			}
		}
	}

	return
}

// IsExistsFn checks whether the path exists in the config file
func IsExistsFn(path string) bool {
	return nil != viper.Get(path)
}

// Get gets the interface{} value from config file by the key
func Get(path string) interface{} {
	return viper.Get(path)
}

// GetInt gets the int value from config file by the key
func GetInt(path string) int {
	return viper.GetInt(path)
}

// GetInt32 gets the int32 value from config file by the key
func GetInt32(path string) int32 {
	return viper.GetInt32(path)
}

// GetInt64 gets the int64 value from config file by the key
func GetInt64(path string) int64 {
	return viper.GetInt64(path)
}

// GetFloat64 gets the float64 value from config file by the key
func GetFloat64(path string) float64 {
	return viper.GetFloat64(path)
}

// GetString gets the string value from config file by the key
func GetString(path string) string {
	return viper.GetString(path)
}

// Set sets the value into path
func Set(path string, value interface{}) {
	viper.Set(path, value)
}

// GetBool gets the boolean value from config file by the key
func GetBool(path string) bool {
	return viper.GetBool(path)
}

// GetUint gets the uint value from config file by the key
func GetUint(path string) uint {
	return viper.GetUint(path)
}

// GetUint32 gets the uint32 value from config file by the key
func GetUint32(path string) uint32 {
	return viper.GetUint32(path)
}

// GetUint64 gets the uint64 value from config file by the key
func GetUint64(path string) uint64 {
	return viper.GetUint64(path)
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal)
}
