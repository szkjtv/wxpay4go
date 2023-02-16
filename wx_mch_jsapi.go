package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MchJsApi struct {
	WxPayMch
}

func NewMchJsApi(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MchJsApi {
	ent := &MchJsApi{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MchJsApi) CreateOrder(data JsApiOrderData) (WxAppPayParam, error) {
	return ent.CreateOrderJsApi(data)
}
