package wxpay

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type PayApi struct {
	Config PayConfig
}

func NewApi(config PayConfig) *PayApi {
	v := new(PayApi)
	v.Config = config
	return v
}

// 统一下单，PayUnifiedOrder中out_trade_no、body、total_fee、trade_type必填
// appid、mchid、spbill_create_ip、nonce_str不需要填入
func (api PayApi) UnifiedOrder(inputObj PayUnifiedOrder, timeout time.Duration) (map[string]interface{}, error) {
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	var (
		ok       bool
		freeType string
	)
	if _, ok := inputObj.GetOutTradeNo(); !ok {
		return nil, errors.New("not set out_trade_no!")
	}
	if _, ok := inputObj.GetBody(); !ok {
		return nil, errors.New("not set body!")
	}
	if _, ok := inputObj.GetTotalFee(); !ok {
		return nil, errors.New("not set total_fee!")
	}
	if freeType, ok = inputObj.GetFeeType(); !ok {
		return nil, errors.New("not set trade_type!")
	}
	if freeType == "JSAPI" {
		if _, ok := inputObj.GetOpenId(); !ok {
			return nil, errors.New("then free_type eq JSAPI, not set openid!")
		}
	}
	if freeType == "NATIVE" {
		if _, ok := inputObj.GetProductId(); !ok {
			return nil, errors.New("then free_type eq NATIVE, not set product_id!")
		}
	}
	if _, ok := inputObj.GetNotifyUrl(); !ok {
		inputObj.SetNotifyUrl(api.Config.NotifyUrl)
	}
	inputObj.SetAppId(api.Config.AppId)
	inputObj.SetMchId(api.Config.MchId)
	//inputObj.SetSpbillCreateIp()
	//inputObj.SetNonceStr()
	inputObj.SetSign(api.Config.Key)
	xmlString, err := inputObj.ToXml()
	if err != nil {
		return nil, err
	}
	data, err := api.postXmlCurl(xmlString, url, false, timeout)
	if err != nil {
		return nil, err
	}
	result := NewPayResults()
	if err := result.Init(data); err != nil {
		return nil, err
	}
	return nil, nil
}

func (api PayApi) postXmlCurl(xmlString, url string, useCert bool, timeout time.Duration) ([]byte, error) {
	var httpTransport *http.Transport
	if api.Config.SSLCertPath != "" {
		cert, err := tls.LoadX509KeyPair(api.Config.SSLCertPath, api.Config.SSLKeyPath)
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		httpTransport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	buff := &bytes.Buffer{}
	buff.WriteString(xmlString)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}
	cli := &http.Client{Transport: httpTransport}
	cli.Timeout = timeout
	// TODO 代理
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	return bts, err
}
