package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
)

type MchH5 struct {
	WxPayMch
}

func NewMchH5(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *MchH5 {
	ent := &MchH5{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付统一下单
func (ent *MchH5) CreateOrder(data WxH5OrderData) (string, error) {
	return ent.CreateOrderH5(data)
}
