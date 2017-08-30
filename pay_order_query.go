package wxpay

type PayOrderQuery interface {
	SetAppId(string)
	GetAppId() (string, bool)
	SetMchId(string)
	GetMchId() (string, bool)
	SetTransactionId(string)
	GetTransactionId() (string, bool)
	SetOutTradeNo(string)
	GetOutTradeNo() (string, bool)
}

func NewPayOrderQuery() PayOrderQuery {
	return NewPayDataBase()
}
