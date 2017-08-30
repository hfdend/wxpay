package wxpay

// 接口调用结果类
type PayResults interface {
	CheckSign(key string) error
	Init(xmlData []byte) error
	FromArray(ary map[string]interface{})
	//InitFromArray(ary map[string]interface{}, noCheckSign bool) (WxPayResults, error)
}

func NewPayResults() PayResults {
	return NewPayDataBase()
}
