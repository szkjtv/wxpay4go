package wxpay

import (
	"crypto/rsa"
	"fmt"
)

//返回给小程序的微信支付数据
type RawAppPayParam struct{
	Appid string
	PartnerId string
	Prepay_id string
	Package string
	NonceStr string
	TimeStamp string
	PaySign string
}

//使用商户私钥对签名串进行SHA256 with RSA签名，并对签名结果进行Base64编码得到签名值
func (ent *RawAppPayParam)GenPaySignV3(mch_pem_key *rsa.PrivateKey) (string, error) {
	SignatureMessageFormat := "%s\n%s\n%s\n%s\n"
	message := fmt.Sprintf(SignatureMessageFormat, ent.Appid, ent.TimeStamp, ent.NonceStr, ent.Prepay_id)
	signatureResult, err := SignSHA256WithRSA(mch_pem_key, message)
	if err != nil {
		return "", err
	}
	return signatureResult, nil
}
