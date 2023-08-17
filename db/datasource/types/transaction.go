package types

import "git.multiverse.io/eventkit/kit/contexts"

type TransactionContexts struct {
	contexts.TransactionContexts
	RoundImages *RoundRecordImage
}
