package mesh

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"github.com/valyala/fasthttp"
)

type meshClient struct {
	opts client.Options
}

func (m meshClient) String() string {
	return fmt.Sprintf(`%++v`, m.opts)
}

var (
	srcDirName       = "src/"
	srcDirNameLength = len(srcDirName)
)

var (
	maxReadTimeoutMilliseconds  = int64(80 * 1000)
	maxWriteTimeoutMilliseconds = int64(80 * 1000)
	once                        sync.Once
	clientWithConnectTimeout    *fasthttp.Client
)

func maskFilePath(filePath string) string {
	if "" != filePath {
		idx := strings.Index(filePath, srcDirName)
		if -1 != idx {
			return filePath[idx+srcDirNameLength:]
		}
	}

	return filePath
}

// NewMeshClient creates a new client.Client with default config,
// and It's easy add or override one or more additional parameters to the client.Client
func NewMeshClient(options ...client.Option) client.Client {
	opts := client.NewOptions()

	for _, o := range options {
		o(&opts)
	}

	mc := &meshClient{
		opts: opts,
	}
	if log.IsEnable(log.DebugLevel) {
		_, file, line, _ := runtime.Caller(1)
		log.Debugsf("new mesh client:[%++v], caller info={file=[%s], line=[%d]}", mc, maskFilePath(file), line)
	}
	return client.Client(mc)
}

func setValueIfValueNotEmpty(m map[string]string, key, value string) {
	if "" != value {
		m[key] = value
	}
}

func setKeyValue(m map[string]string, key, value string) {
	m[key] = value
}

func createClientIfNecessary() {
	once.Do(func() {
		clientWithConnectTimeout = &fasthttp.Client{
			// Default value of MaxConnsPerHost is 512, increase to 16384
			MaxConnsPerHost: 16384,
			// Disable idempotent calls attempts when remote call abnormal.
			//
			// default max attempts is 5 times.
			MaxIdemponentCallAttempts: 1,
		}
		log.Infosf("Create http client with read timeout(milliseconds)=[%d], write timeout(milliseconds)=[%d]", maxReadTimeoutMilliseconds, maxWriteTimeoutMilliseconds)
		clientWithConnectTimeout.ReadTimeout = time.Duration(maxReadTimeoutMilliseconds) * time.Millisecond
		clientWithConnectTimeout.WriteTimeout = time.Duration(maxWriteTimeoutMilliseconds) * time.Millisecond
	})
}

func injectTopicAttributes(message *msg.Message, opts *client.RequestOptions) error {
	topicAttributes := make(map[string]string)

	switch opts.TopicType {
	case constant.TopicTypeBusiness, constant.TopicTypeOPS:
		{
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationORG, opts.Org)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationWorkspace, opts.Wks)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationEnvironment, opts.Env)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationSU, opts.Su)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationVersion, opts.Version)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicID, opts.EventID)
		}
	case constant.TopicTypeHeartbeat, constant.TopicTypeError,
		constant.TopicTypeAlert, constant.TopicTypeLog,
		constant.TopicTypeMetrics:
		{
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationORG, opts.Org)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationWorkspace, opts.Wks)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationEnvironment, opts.Env)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicID, opts.EventID)
		}
	case constant.TopicTypeDXC, constant.TopicTypeDTS:
		{
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationORG, opts.Org)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationWorkspace, opts.Wks)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationEnvironment, opts.Env)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationSU, opts.Su)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationDCN, opts.Su)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationNodeID, opts.NodeID)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicDestinationInstanceID, opts.InstanceID)
			setValueIfValueNotEmpty(topicAttributes, constant.TopicID, opts.EventID)
		}
	default:
		return errors.Errorf(constant.SystemInternalError, "unsupported topic type[%s]", opts.TopicType)
	}

	topicAttributes[constant.TopicType] = opts.TopicType
	message.TopicAttribute = topicAttributes

	return nil
}

func injectRequestHeader(message *msg.Message, opts *client.RequestOptions) {
	opts.HeaderLock.RLock()
	if nil != opts.Header {
		for k, v := range opts.Header {
			message.SetAppProperty(k, v)
		}
	}
	opts.HeaderLock.RUnlock()

	if len(opts.SessionName) > 0 {
		message.SessionName = opts.SessionName
	}

	message.SetAppProperty(constant.To3, strconv.FormatInt(int64(opts.Timeout.Seconds()*1000), 10))

	// set persistent delivery mode
	if opts.IsPersistentDeliveryMode {
		message.SetAppProperty(constant.DeliveryMode, constant.Enable)
	}

	// set semi-sync call
	if opts.IsSemiSyncCall {
		message.SetAppProperty(constant.TxnIsSemiSyncCall, constant.Enable)
	}

	// set enable DMQ eligible
	if opts.IsDMQEligible {
		message.SetAppProperty(constant.DmqEligible, constant.Enable)
	}

	// set the pass through header
	if len(opts.PassThroughHeaderKeyList) > 0 {
		message.SetAppProperty(constant.PassThroughHeaderKeyListKey, strings.Join(opts.PassThroughHeaderKeyList, ","))
		if len(opts.OriginalHeader) > 0 {
			for _, key := range opts.PassThroughHeaderKeyList {
				if v, ok1 := opts.OriginalHeader[key]; ok1 {
					if _, ok2 := message.GetAppProperty(key); !ok2 {
						message.SetAppProperty(key, v)
					}
				}
			}
		}
	}
}

func injectResponseHeader(message *msg.Message, opts *client.ResponseOptions) {
	if len(opts.SessionName) > 0 {
		message.SessionName = opts.SessionName
	}
	if nil != opts.Header {
		for k, v := range opts.Header {
			message.SetAppProperty(k, v)
		}
	}
}

// Tuple2 A tuple that holds two values.
type Tuple2 struct {
	F0 interface{}
	F1 interface{}
}

func httpRequest(ctx context.Context, request client.Request, requestMessage *msg.Message, sync bool) (responseMessage *msg.Message, err error) {
	// using http to call downstream service
	createClientIfNecessary()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.Set(constant.HTTPContentTypeKey, request.RequestOptions().ContentType)
	req.Header.SetMethod(request.RequestOptions().HTTPMethod)
	req.Header.DisableNormalizing()
	// inject header
	requestMessage.RangeAppProps(func(k string, v string) {
		req.Header.Add(k, v)
	})

	req.SetRequestURI(request.RequestOptions().Address)
	req.SetBody(requestMessage.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if sync {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, request.RequestOptions().Timeout)
		defer cancel()

		deadline, _ := ctx.Deadline()
		err = clientWithConnectTimeout.DoDeadline(req, resp, deadline)
	} else {
		err = clientWithConnectTimeout.Do(req, resp)
	}

	if nil != err {
		if fasthttp.ErrTimeout == err {
			return nil, errors.Wrap(constant.SystemRemoteCallTimeout, err, 0)
		}

		if fasthttp.ErrConnectionClosed == err {
			return nil, errors.Errorf(constant.SystemErrConnectionClosed, "Connection closed, url=[%s]", request.RequestOptions().Address)
		}

		switch err.(type) {
		case *net.OpError:
			{
				netOpError := err.(*net.OpError)
				switch netOpError.Err.(type) {
				case *os.SyscallError:
					{
						syscallError := netOpError.Err.(*os.SyscallError)
						if errno, ok := syscallError.Err.(syscall.Errno); ok {
							switch errno {
							case syscall.ECONNREFUSED:
								{
									return nil, errors.Wrap(constant.SystemErrConnectionRefused, err, 0)
								}
							case syscall.ECONNRESET:
								{
									return nil, errors.Wrap(constant.SystemErrConnectionReset, err, 0)
								}
							case syscall.ECONNABORTED:
								{
									return nil, errors.Wrap(constant.SystemErrConnectionAborted, err, 0)
								}
							default:
								return nil, errors.Wrap(constant.SystemInternalError, err, 0)
							}
						}
					}
				default:
					return nil, errors.Wrap(constant.SystemInternalError, err, 0)
				}
			}
		default:
			return nil, errors.Wrap(constant.SystemInternalError, err, 0)
		}
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.Errorf(strconv.Itoa(resp.StatusCode()), "request StatusCode != 200, url=%s, statusCode=%v",
			request.RequestOptions().Address, resp.StatusCode())
	}

	if sync {
		// copy body to target
		resData := resp.Body()
		body := make([]byte, len(resData))
		copy(body, resData)

		replyHeader := make(map[string]string)

		resp.Header.VisitAll(func(key, value []byte) {
			if len(key) > 0 {
				replyHeader[string(key)] = string(value)
			}
		})
		msg := &msg.Message{
			Body: body,
		}
		msg.SetAppProps(replyHeader)
		return msg, nil
	}

	return nil, nil
}

func (m *meshClient) ReplySemiSyncCall(ctx context.Context, response client.Response) error {
	responseOptions := response.ResponseOptions()

	// encode response
	responseBody, err := responseOptions.Codec.Encoder().Encode(response.Body())
	if nil != err {
		return errors.New(constant.UpstreamServiceMessageEncodeError, err)
	}

	responseMessage := &msg.Message{
		Body: responseBody,
	}

	// inject response header
	injectResponseHeader(responseMessage, responseOptions)

	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		handlerContexts.Inject(true, responseMessage)
	}

	responseMessage.DeleteProperty(constant.TxnIsSemiSyncCall)
	responseMessage.SetAppProperty(constant.RrReplyTo, responseOptions.ReplyToAddress)

	err = callback.ReplySemiSyncCall(responseMessage)

	if !util.IsNil(err) {
		// for enable debug stack
		var retErr *errors.Error
		switch err.(type) {
		case *errors.Error:
			e := err.(*errors.Error)
			log.Errorf(ctx, "Reply semi sync call[reply to address:%s], error:[%s]", responseOptions.ReplyToAddress, e.Error())
			retErr = e
		default:
			retErr = errors.Wrap(constant.SystemInternalError, fmt.Sprintf("Reply semi sync call[reply to address:%s], error:[%++v]", responseOptions.ReplyToAddress, err), 0)
		}

		log.Errorf(ctx, "Reply semi sync call[reply to address:%s], error:[%s], call stack:[%s]", responseOptions.ReplyToAddress, retErr.Error(), retErr.ErrorStack())
		return retErr
	}

	return nil
}

func (m *meshClient) SyncCall(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (res client.ResponseMeta, oerr error) {
	callOpts := m.opts.CallOptions
	requestOptions := request.RequestOptions()

	for _, opt := range opts {
		opt(&callOpts)
	}

	defer func() {
		if oerr != nil && nil == res {
			errorCode := ""
			errorMessage := ""
			switch oerr.(type) {
			case *errors.Error:
				{
					e := oerr.(*errors.Error)
					errorCode = e.ErrorCode
					errorMessage = e.Error()
				}
			default:
				{
					errorCode = constant.SystemInternalError
					errorMessage = fmt.Sprintf("%++v", oerr)
				}
			}
			res = NewMeshResponseMeta(nil, map[string]string{
				constant.ReturnErrorCode: errorCode,
				constant.ReturnErrorMsg:  errorMessage,
			})
		}

		// do Wrappers
		for _, wrapper := range callOpts.CallWrappers {
			wrapperName := fmt.Sprintf("%s", wrapper)
			if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallWrapper(wrapperName) {
				log.Debugf(ctx, "SyncCall,Skipping the `After` of remote call wrapper:%s", wrapperName)
			} else {
				tmpCtx, wErr := wrapper.After(ctx, request, res, requestOptions)
				if nil != wErr {
					log.Errorf(ctx, "SyncCall,Failed do wrapper %s , error:%++v", wrapper, wErr)
				}
				ctx = tmpCtx
			}
		}
	}()
	d, ok := ctx.Deadline()

	if !ok {
		if requestOptions.MaxWaitingTime.Seconds() > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, requestOptions.MaxWaitingTime)
			defer cancel()
		}
	} else {
		opt := WithTimeout(d.Sub(time.Now()))
		opt(requestOptions)
	}

	// do Wrappers
	for _, wrapper := range callOpts.CallWrappers {
		wrapperName := fmt.Sprintf("%s", wrapper)
		if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallWrapper(wrapperName) {
			log.Debugf(ctx, "SyncCall,Skipping the `Before` of remote call wrapper:%s", wrapperName)
		} else {
			tmpCtx, wErr := wrapper.Before(ctx, request, requestOptions)
			if nil != wErr {
				return nil, wErr
			}
			ctx = tmpCtx
		}
	}

	// encode request
	requestBody, err := request.Codec().Encoder().Encode(request.Body())
	if nil != err {
		return nil, errors.New(constant.DownstreamServiceMessageEncodeError, err)
	}

	requestMessage := &msg.Message{
		Body: requestBody,
	}

	// inject topic attributes
	if !request.RequestOptions().HTTPCall {
		if err = injectTopicAttributes(requestMessage, requestOptions); nil != err {
			return nil, err
		}
	}
	// inject header
	injectRequestHeader(requestMessage, requestOptions)

	// inject context into appProps
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		handlerContexts.Inject(request.RequestOptions().DeleteTransactionPropagationInformation, requestMessage)
	}

	if !requestOptions.IsSemiSyncCall {
		requestMessage.DeleteProperty(constant.TxnIsSemiSyncCall)
		requestMessage.DeleteProperty(constant.RrReplyTo)
	}

	call := func(i int) *Tuple2 {
		var responseMessage *msg.Message
		var indexOfInterceptorsExecuted int

		t, err := requestOptions.Backoff(ctx, request, i)
		if nil != err {
			return &Tuple2{
				F0: nil,
				F1: errors.Errorf(constant.SystemInternalError, "backoff error:%v", err),
			}
		}
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// do preHandle of interceptors
		for i := 0; i < len(callOpts.CallInterceptors); i++ {
			interceptor := callOpts.CallInterceptors[i]
			interceptorName := fmt.Sprintf("%s", interceptor)
			if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallInterceptor(interceptorName) {
				log.Debugf(ctx, "SyncCall,Skipping the `PreHandle` of remote call interceptor:%s", interceptorName)
			} else {
				if ierr := interceptor.PreHandle(ctx, requestMessage); nil != ierr {
					return &Tuple2{
						F0: nil,
						F1: ierr,
					}
				}
			}

			indexOfInterceptorsExecuted++
		}
		requestMessage.DeleteProperty(constant.TxnIsLocalCall)
		if requestOptions.IsLocalCall {
			requestMessage.SetAppProperty(constant.TxnIsLocalCall, "1")
			if requestOptions.EnableLogging {
				log.Debugf(ctx, "call local service in SyncCall, request:[%s], request Options:[%++v]", requestMessage, request.RequestOptions())
			} else {
				log.Debugf(ctx, "call local service in SyncCall, request topic attribute:[%s], request Options:[%++v]", requestMessage.TopicAttribute, request.RequestOptions())
			}
			responseMessage, err = callOpts.CallbackExecutor.Handle(context.Background(), requestMessage)
		} else if request.RequestOptions().HTTPCall {
			if requestOptions.EnableLogging {
				log.Debugf(ctx, "client using http to send message in SyncCall, request:[%s], request Options:[%++v]", requestMessage, request.RequestOptions())
			} else {
				log.Debugf(ctx, "client using http to send message in SyncCall, request topic attribute:[%s], request Options:[%++v]", requestMessage.TopicAttribute, request.RequestOptions())
			}
			responseMessage, err = httpRequest(ctx, request, requestMessage, true)
		} else {
			if requestOptions.EnableLogging {
				log.Debugf(ctx, "client using mesh to send message in SyncCall, request:[%s], request Options:[%++v]", requestMessage, request.RequestOptions())
			} else {
				log.Debugf(ctx, "client using mesh to send message in SyncCall, request topic attribute:[%s], request Options:[%++v]", requestMessage.TopicAttribute, request.RequestOptions())
			}
			responseMessage, err = callback.SyncCall(requestMessage, request.RequestOptions().Timeout)
		}
		if err == nil {
			if requestOptions.EnableLogging {
				log.Debugf(ctx, "client receive message in SyncCall, header:[%s] body:[%s]",
					responseMessage.AppPropsToString(), string(responseMessage.Body))
			} else {
				log.Debugf(ctx, "client receive message in SyncCall, topic attribute:[%s]", responseMessage.TopicAttribute)
			}
		}

		if nil != err {
			return &Tuple2{
				F0: nil,
				F1: err,
			}
		}

		for i := indexOfInterceptorsExecuted - 1; i >= 0; i-- {
			interceptor := callOpts.CallInterceptors[i]
			interceptorName := fmt.Sprintf("%s", interceptor)
			if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallInterceptor(interceptorName) {
				log.Debugf(ctx, "Skipping the `PostHandle` of remote call interceptor:%s", interceptorName)
			} else {
				if ierr := interceptor.PostHandle(ctx, requestMessage, responseMessage); nil != ierr {
					return &Tuple2{
						F0: nil,
						F1: ierr,
					}
				}
			}
		}
		if err := request.Codec().Decoder().Decode(responseMessage.Body, response); nil != err {
			return &Tuple2{
				F0: nil,
				F1: errors.New(constant.DownstreamServiceMessageDecodeError, err),
			}
		}

		return &Tuple2{
			F0: responseMessage,
			F1: nil,
		}
	}
	retries := requestOptions.MaxRetryTimes

	ch := make(chan *Tuple2, retries+1)
	var e error

	for i := 0; i <= retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return nil, errors.Errorf(constant.SystemRemoteCallTimeout, "call timeout: %v", ctx.Err())
		case tuple2 := <-ch:
			var retErr *errors.Error
			if util.IsNil(tuple2.F1) {
				resMsg := tuple2.F0.(*msg.Message)
				return NewMeshResponseMeta(resMsg.Body, resMsg.GetAppProps()), nil
			}

			// for enable debug stack
			switch tuple2.F1.(type) {
			case *errors.Error:
				e := tuple2.F1.(*errors.Error)
				log.Errorf(ctx, "SyncCall request attributes:%s, error:[%s]", util.MapToString(requestMessage.TopicAttribute), e.Error())
				retErr = e
			case error:
				retErr = errors.Wrap(constant.SystemInternalError, fmt.Sprintf("SyncCall request attributes:%s, error:[%++v]", util.MapToString(requestMessage.TopicAttribute), tuple2.F1.(error)), 0)
			default:
				retErr = errors.Errorf(constant.SystemInternalError, "SyncCall request attributes:%s, info:%++v", util.MapToString(requestMessage.TopicAttribute), e)
			}

			log.Errorf(ctx, "Failed to SyncCall [SEQ=%d]request attributes:%s, error:[%s], call stack:[%s]", i+1, util.MapToString(requestMessage.TopicAttribute), retErr.Error(), retErr.ErrorStack())
			retry, rerr := requestOptions.Retry(ctx, request, i, retErr)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, retErr
			}

			e = retErr
		}
	}

	return nil, e
}

func (m *meshClient) AsyncCall(ctx context.Context, request client.Request, opts ...client.CallOption) error {
	callOpts := m.opts.CallOptions
	requestOptions := request.RequestOptions()

	for _, opt := range opts {
		opt(&callOpts)
	}

	defer func() {
		// do wrappers
		for _, wrapper := range callOpts.CallWrappers {
			wrapperName := fmt.Sprintf("%s", wrapper)
			if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallWrapper(wrapperName) {
				log.Debugf(ctx, "AsyncCall,Skipping the `After` of remote call wrapper:%s", wrapperName)
			} else {
				tmpCtx, wErr := wrapper.After(ctx, request, nil, requestOptions)
				if nil != wErr {
					log.Errorf(ctx, "Failed do wrapper %s , error:%++v", wrapper, wErr)
				}
				ctx = tmpCtx
			}
		}
	}()

	// do wrappers
	for _, wrapper := range callOpts.CallWrappers {
		wrapperName := fmt.Sprintf("%s", wrapper)
		if nil != requestOptions.ServiceConfig && requestOptions.ServiceConfig.IsMarkedAsSkippedRemoteCallWrapper(wrapperName) {
			log.Debugf(ctx, "AsyncCall,Skipping the `Before` of remote call wrapper:%s", wrapperName)
		} else {
			tmpCtx, wErr := wrapper.Before(ctx, request, requestOptions)
			if nil != wErr {
				return wErr
			}
			ctx = tmpCtx
		}
	}

	// encode request
	requestBody, err := request.Codec().Encoder().Encode(request.Body())
	if nil != err {
		return errors.New(constant.DownstreamServiceMessageEncodeError, err)
	}

	requestMessage := &msg.Message{
		Body: requestBody,
	}

	// inject topic attributes
	if !request.RequestOptions().HTTPCall {
		if err = injectTopicAttributes(requestMessage, requestOptions); nil != err {
			return err
		}
	}

	// inject header
	injectRequestHeader(requestMessage, requestOptions)

	// inject context into appProps
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		handlerContexts.Inject(request.RequestOptions().DeleteTransactionPropagationInformation, requestMessage)
	}
	requestMessage.DeleteProperty(constant.TxnIsLocalCall)
	if request.RequestOptions().HTTPCall {
		if requestOptions.EnableLogging {
			log.Debugf(ctx, "client using http to send message in AsyncCall, request:[%s]", requestMessage)
		} else {
			log.Debugf(ctx, "client using http to send message in AsyncCall, request topic attribute:[%s]", util.MapToString(requestMessage.TopicAttribute))
		}
		_, err = httpRequest(ctx, request, requestMessage, false)
	} else {
		if requestOptions.EnableLogging {
			log.Debugf(ctx, "client using mesh to send message in AsyncCall, request:[%s]", requestMessage)
		} else {
			log.Debugf(ctx, "client using mesh to send message in AsyncCall, request topic attribute:[%s]", util.MapToString(requestMessage.TopicAttribute))
		}

		err = callback.Publish(requestMessage)
	}

	if !util.IsNil(err) {
		// for enable debug stack
		var retErr *errors.Error
		switch err.(type) {
		case *errors.Error:
			e := err.(*errors.Error)
			log.Errorf(ctx, "AsyncCall request attributes:%s, error:[%s]", util.MapToString(requestMessage.TopicAttribute), e.Error())
			retErr = e
		default:
			retErr = errors.Wrap(constant.SystemInternalError, fmt.Sprintf("AsyncCall request attributes:%s, error:[%++v]", util.MapToString(requestMessage.TopicAttribute), err), 0)
		}

		log.Errorf(ctx, "Failed to AsyncCall request attributes:%s, error:[%s], call stack:[%s]", util.MapToString(requestMessage.TopicAttribute), retErr.Error(), retErr.ErrorStack())
		return retErr
	}

	return nil
}

func (m *meshClient) Options() client.Options {
	return m.opts
}
