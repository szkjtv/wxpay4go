package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MidH5 struct {
	WxPayMid
}

func NewMidH5(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MidH5 {
	ent := &MidH5{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MidH5) CreateOrder(data WxH5OrderData) (string, error) {
	return ent.CreateOrderH5(data)
}
