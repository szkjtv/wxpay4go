package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MchNative struct {
	WxPayMch
}

func NewMchNative(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MchNative {
	ent := &MchNative{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MchNative) CreateOrder(data NativeOrderData) (string, error) {
	return ent.CreateOrderNative(data)
}
