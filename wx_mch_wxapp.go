package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MchWxApp struct {
	WxPayMch
}

func NewMchWxApp(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MchWxApp {
	ent := &MchWxApp{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MchWxApp) CreateOrder(data JsApiOrderData) (WxAppPayParam, error) {
	return ent.CreateOrderJsApi(data)
}
