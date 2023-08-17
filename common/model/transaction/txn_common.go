package transaction

//TxnEventHeader is the header of transaction callback model
type TxnEventHeader struct {
	Service string `json:"service"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
