package config

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestGLSCache_Equals(t *testing.T) {
	v1 := GLSCache{
		Type:     "Type",
		Addr:     "Addr",
		Password: "Password",
		PoolNum:  10,
		Readonly: true,
	}

	v2 := &GLSCache{
		Type:     "Type",
		Addr:     "Addr",
		Password: "Password",
		PoolNum:  10,
		Readonly: true,
	}

	assert.True(t, v1.Equals(v2))

	v2.Type = "new type"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestService_Equals(t *testing.T) {
	v1 := Service{
		ServiceID:        "ServiceID",
		Org:              "Org",
		Az:               "Az",
		Wks:              "Wks",
		Env:              "Env",
		NodeID:           "NodeID",
		InstanceID:       "InstanceID",
		Su:               "Su",
		GroupSu:          "GroupSu",
		CommonSu:         "CommonSu",
		ResponseTemplate: "ResponseTemplate",
		ResponseAutoParseKeyMapping: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
	}

	v2 := &Service{
		ServiceID:        "ServiceID",
		Org:              "Org",
		Az:               "Az",
		Wks:              "Wks",
		Env:              "Env",
		NodeID:           "NodeID",
		InstanceID:       "InstanceID",
		Su:               "Su",
		GroupSu:          "GroupSu",
		CommonSu:         "CommonSu",
		ResponseTemplate: "ResponseTemplate",
		ResponseAutoParseKeyMapping: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
	}

	assert.True(t, v1.Equals(v2))

	v2.Su = "new su"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestLog_Equals(t *testing.T) {
	v1 := Log{
		LogFile:            "LogFile",
		LogLevel:           "LogLevel",
		LogLevelUnixSocket: "LogLevelUnixSocket",
		LogFileRootPath:    "LogFileRootPath",
		MaxSize:            100,
		MaxDays:            200,
		MaxBackups:         100,
		Console:            false,
		MaskRules:          nil,
	}

	v2 := &Log{
		LogFile:            "LogFile",
		LogLevel:           "LogLevel",
		LogLevelUnixSocket: "LogLevelUnixSocket",
		LogFileRootPath:    "LogFileRootPath",
		MaxSize:            100,
		MaxDays:            200,
		MaxBackups:         100,
		Console:            false,
		MaskRules:          nil,
	}

	assert.True(t, v1.Equals(v2))

	v2.MaxSize = 120
	assert.False(t, v1.Equals(v2))
	assert.False(t, v1.Equals(nil))
}

func TestDb_Equals(t *testing.T) {
	v1 := Db{
		Name:     "Name",
		Type:     "Type",
		Su:       "Su",
		Topics:   []string{"t1", "t2", "t3"},
		Default:  true,
		Addr:     "Addr",
		User:     "User",
		Password: "Password",
		Database: "Database",
		Params:   "Params",
		Debug:    true,
		Pool: struct {
			MaxIdleConns      int `json:"maxIdleConns"`
			MaxOpenConns      int `json:"maxOpenConns"`
			MaxIdleTime 	  int `json:"maxIdleTime"`
			MaxLifeValue      int `json:"maxLifeValue"`
		}{
			MaxIdleConns:      100,
			MaxOpenConns:      200,
			MaxIdleTime: 	   100,
			MaxLifeValue:      500,
		},
	}

	v2 := &Db{
		Name:     "Name",
		Type:     "Type",
		Su:       "Su",
		Topics:   []string{"t1", "t2", "t3"},
		Default:  true,
		Addr:     "Addr",
		User:     "User",
		Password: "Password",
		Database: "Database",
		Params:   "Params",
		Debug:    true,
		Pool: struct {
			MaxIdleConns      int `json:"maxIdleConns"`
			MaxOpenConns      int `json:"maxOpenConns"`
			MaxIdleTime 	  int `json:"maxIdleTime"`
			MaxLifeValue      int `json:"maxLifeValue"`
		}{
			MaxIdleConns:      100,
			MaxOpenConns:      200,
			MaxIdleTime: 	   100,
			MaxLifeValue:      500,
		},
	}

	assert.True(t, v1.Equals(v2))

	v2.Pool.MaxLifeValue = 1
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestTransaction_Equals(t *testing.T) {
	v1 := Transaction{
		CommType:                      "CommType",
		IsPropagator:                  false,
		TryFailedIgnoreCallbackCancel: true,
		SaveHeaders:                   false,
		TimeoutMilliseconds:           10,
		MaxRetryTimes:                 20,
		MaxServiceConsumeMilliseconds: 30,
		PropagatorServices:            []string{"s1", "s2", "s3"},
		PropagatorServicesMap: map[string]bool{
			"k1": true,
			"k2": true,
			"k3": true,
		},
		TransactionServer: TransactionServer{
			Org:             "Org",
			Wks:             "Wks",
			Env:             "Env",
			Su:              "Su",
			NodeID:          "NodeID",
			InstanceID:      "InstanceID",
			TxnBeginEventID: "TxnBeginEventID",
			TxnJoinEventID:  "TxnJoinEventID",
			TxnEndEventID:   "TxnEndEventID",
			AddressURL:      "AddressURL",
			TxnBeginURLPath: "TxnBeginURLPath",
			TxnJoinURLPath:  "TxnJoinURLPath",
			TxnEndURLPath:   "TxnEndURLPath",
		},
		TransactionClient: TransactionClient{
			ConfirmEventID:     "ConfirmEventID",
			CancelEventID:      "CancelEventID",
			ParticipantAddress: "ParticipantAddress",
			ConfirmAddressURL:  "ConfirmAddressURL",
			CancelAddressURL:   "CancelAddressURL",
		},
	}

	v2 := &Transaction{
		CommType:                      "CommType",
		IsPropagator:                  false,
		TryFailedIgnoreCallbackCancel: true,
		SaveHeaders:                   false,
		TimeoutMilliseconds:           10,
		MaxRetryTimes:                 20,
		MaxServiceConsumeMilliseconds: 30,
		PropagatorServices:            []string{"s1", "s2", "s3"},
		PropagatorServicesMap: map[string]bool{
			"k1": true,
			"k2": true,
			"k3": true,
		},
		TransactionServer: TransactionServer{
			Org:             "Org",
			Wks:             "Wks",
			Env:             "Env",
			Su:              "Su",
			NodeID:          "NodeID",
			InstanceID:      "InstanceID",
			TxnBeginEventID: "TxnBeginEventID",
			TxnJoinEventID:  "TxnJoinEventID",
			TxnEndEventID:   "TxnEndEventID",
			AddressURL:      "AddressURL",
			TxnBeginURLPath: "TxnBeginURLPath",
			TxnJoinURLPath:  "TxnJoinURLPath",
			TxnEndURLPath:   "TxnEndURLPath",
		},
		TransactionClient: TransactionClient{
			ConfirmEventID:     "ConfirmEventID",
			CancelEventID:      "CancelEventID",
			ParticipantAddress: "ParticipantAddress",
			ConfirmAddressURL:  "ConfirmAddressURL",
			CancelAddressURL:   "CancelAddressURL",
		},
	}

	assert.True(t, v1.Equals(v2))

	v2.TransactionClient.ParticipantAddress = "new participant address"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestHeartbeat_Equals(t *testing.T) {
	v1 := Heartbeat{
		TopicName:       "TopicName",
		IntervalSeconds: 100,
	}

	v2 := &Heartbeat{
		TopicName:       "TopicName",
		IntervalSeconds: 100,
	}

	assert.True(t, v1.Equals(v2))

	v2.IntervalSeconds = 1
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestAlert_Equals(t *testing.T) {
	v1 := Alert{
		TopicName: "topic name",
	}

	v2 := &Alert{
		TopicName: "topic name",
	}
	assert.True(t, v1.Equals(v2))

	v2.TopicName = "new topic name"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestApm_Equals(t *testing.T) {
	v1 := Apm{
		Enable:                              true,
		PrintEmptyTraceIdRecordAtClientSide: false,
		Version:                             "v1",
		RootPath:                            "RootPath",
		FileRows:                            100,
	}

	v2 := &Apm{
		Enable:                              true,
		PrintEmptyTraceIdRecordAtClientSide: false,
		Version:                             "v1",
		RootPath:                            "RootPath",
		FileRows:                            100,
	}

	assert.True(t, v1.Equals(v2))

	v2.Enable = false
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestDeployment_Equals(t *testing.T) {
	v1 := Deployment{
		EnableSecure: true,
	}

	v2 := &Deployment{
		EnableSecure: true,
	}

	assert.True(t, v1.Equals(v2))

	v2.EnableSecure = false
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestAddressing_Equals(t *testing.T) {
	v1 := Addressing{
		Enable:               false,
		SyncConfigWithServer: false,
		TopicIDOfServer:      "TopicIDOfServer",
		TopicSuTitle:         "TopicSuTitle",
		Cache: GLSCache{
			Type:     "Type",
			Addr:     "Addr",
			Password: "Password",
			PoolNum:  10,
			Readonly: true,
		},
	}
	v2 := &Addressing{
		Enable:               false,
		SyncConfigWithServer: false,
		TopicIDOfServer:      "TopicIDOfServer",
		TopicSuTitle:         "TopicSuTitle",
		Cache: GLSCache{
			Type:     "Type",
			Addr:     "Addr",
			Password: "Password",
			PoolNum:  10,
			Readonly: true,
		},
	}

	assert.True(t, v1.Equals(v2))

	v2.TopicSuTitle = "new TopicSuTitle"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestDownstream_Equals(t *testing.T) {
	v1 := Downstream{
		EventType:                        "EventType",
		EventID:                          "EventID",
		Version:                          "Version",
		TimeoutMilliseconds:              10,
		RetryWaitingMilliseconds:         20,
		MaxWaitingTimeMilliseconds:       30,
		MaxRetryTimes:                    40,
		DeleteTransactionPropagationInfo: true,
		ProtoType:                        "ProtoType",
		HTTPAddress:                      "HTTPAddress",
		HTTPMethod:                       "HTTPMethod",
		HTTPContextType:                  "HTTPContextType",
		ResponseAutoParseKeyMapping: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
	}

	v2 := &Downstream{
		EventType:                        "EventType",
		EventID:                          "EventID",
		Version:                          "Version",
		TimeoutMilliseconds:              10,
		RetryWaitingMilliseconds:         20,
		MaxWaitingTimeMilliseconds:       30,
		MaxRetryTimes:                    40,
		DeleteTransactionPropagationInfo: true,
		ProtoType:                        "ProtoType",
		HTTPAddress:                      "HTTPAddress",
		HTTPMethod:                       "HTTPMethod",
		HTTPContextType:                  "HTTPContextType",
		ResponseAutoParseKeyMapping: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
	}

	assert.True(t, v1.Equals(v2))

	v2.Version = "new version"
	assert.False(t, v1.Equals(v2))

	assert.False(t, v1.Equals(nil))
}

func TestService_IsMarkedAsSkippedHandleInterceptor(t *testing.T) {
	service := Service{
		SkipHandleInterceptorsMap: map[string]bool{"apm": true, "trace": true},
	}

	assert.True(t, service.IsMarkedAsSkippedHandleInterceptor("apm"))
	assert.False(t, service.IsMarkedAsSkippedHandleInterceptor("does_not_exist"))
	assert.True(t, service.IsMarkedAsSkippedHandleInterceptor("TRace"))
}

func TestService_IsMarkedAsSkippedRemoteCallInterceptor(t *testing.T) {
	service := Service{
		SkipRemoteCallInterceptorsMap: map[string]bool{"apm": true, "trace": true},
	}
	assert.True(t, service.IsMarkedAsSkippedRemoteCallInterceptor("apm"))
	assert.False(t, service.IsMarkedAsSkippedRemoteCallInterceptor("does_not_exist"))
	assert.True(t, service.IsMarkedAsSkippedRemoteCallInterceptor("TRace"))
}

func TestService_IsMarkedAsSkippedRemoteCallWrapper(t *testing.T) {
	service := Service{
		SkipRemoteCallWrappersMap: map[string]bool{"apm": true, "trace": true},
	}
	assert.True(t, service.IsMarkedAsSkippedRemoteCallWrapper("apm"))
	assert.False(t, service.IsMarkedAsSkippedRemoteCallWrapper("does_not_exist"))
	assert.True(t, service.IsMarkedAsSkippedRemoteCallWrapper("TRace"))
}
