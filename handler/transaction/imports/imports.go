package imports

import (
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/router"
	"git.multiverse.io/eventkit/kit/handler/transaction/callback"
	"git.multiverse.io/eventkit/kit/log"
)

// HasEnabledTransactionSupport is used to mark whether the service already supports transactions
var HasEnabledTransactionSupport = false

// EnableTransactionSupports initialize routers for Transaction
func EnableTransactionSupports(handlerRouter *router.HandlerRouter, transactionClient *config.TransactionClient) (err error) {
	log.Infosf("start register confirm/cancel handler...")

	confirmOptions := []router.Option{
		router.Method("CallbackConfirm"),
		router.WithCodec(codec.BuildTextCodec()),
		router.WithInterceptors(),
		router.DisableValidation(),
	}

	if len(transactionClient.ConfirmAddressURL) > 0 {
		confirmOptions = append(confirmOptions, router.HandlePost(transactionClient.ConfirmAddressURL))
	}

	handlerRouter.Router(transactionClient.ConfirmEventID,
		&callback.TransactionCallbackHandler{},
		confirmOptions...,
	)

	cancelOptions := []router.Option{
		router.Method("CallbackCancel"),
		router.WithCodec(codec.BuildTextCodec()),
		router.WithInterceptors(),
		router.DisableValidation(),
	}
	if len(transactionClient.CancelAddressURL) > 0 {
		cancelOptions = append(cancelOptions, router.HandlePost(transactionClient.CancelAddressURL))
	}

	handlerRouter.Router(transactionClient.CancelEventID,
		&callback.TransactionCallbackHandler{},
		cancelOptions...,
	)

	// mark has executed EnableTransactionSupports
	HasEnabledTransactionSupport = true

	log.Infosf("start register confirm/cancel handler...")
	return
}
