package wxpay

import (
	"crypto/rsa"
	"fmt"
)

//返回给小程序的微信支付数据
type WxAppPayParam struct{
	Appid string
	TimeStamp string
	NonceStr string
	Prepay_id string
	Package string
	SignType string
	PaySign string
}

//使用商户私钥对签名串进行SHA256 with RSA签名，并对签名结果进行Base64编码得到签名值
func (ent *WxAppPayParam)GenPaySignV3(mch_pem_key *rsa.PrivateKey) (string, error) {
	SignatureMessageFormat := "%s\n%s\n%s\n%s\n"
	message := fmt.Sprintf(SignatureMessageFormat, ent.Appid, ent.TimeStamp, ent.NonceStr, ent.Package)
	signatureResult, err := SignSHA256WithRSA(mch_pem_key, message)
	if err != nil {
		return "", err
	}
	return signatureResult, nil
}
