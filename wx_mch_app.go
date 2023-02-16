package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MchApp struct {
	WxPayMch
}

func NewMchApp(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MchApp {
	ent := &MchApp{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MchApp) CreateOrder(data OsAppOrderData) (RawAppPayParam, error) {
	return ent.CreateOrderApp(data)
}
