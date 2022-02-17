package api

const (
	serverName = ""
	rpcPath    = "/_tchain_rpc"
)

type Sender interface {
	SendTransaction(TransactionReq) (TransactionResp, error)
	SendIsAlive() error
	SendBlock() error
}

type Receiver interface {
	HandleTransaction(TransactionReq, *TransactionResp) error
	HandleIsAlive(Empty, *Empty) error
	HandleBlock(BlockReq, *OpStatus) error
}

type (
	Empty struct{}

	OpStatus struct {
		Status bool
		Msg    string
	}
)

type BlockReq struct {
	Block
}

type (
	TransactionReq struct {
		Transaction
	}

	TransactionResp struct {
		Status bool
		Msg    string
	}
)
