package wxpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

//生成32个字节的请求随机串
const NonceSymbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func GenerateNonce() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	symbolsByteLength := byte(len(NonceSymbols))
	for i, b := range bytes {
		bytes[i] = NonceSymbols[b%symbolsByteLength]
	}
	return string(bytes), nil
}

//使用商户私钥对待签名串进行SHA256 with RSA签名，并对签名结果进行Base64编码
func SignSHA256WithRSA(privateKey *rsa.PrivateKey, source string) (signature string, err error) {
	h := crypto.Hash.New(crypto.SHA256)
	_, err = h.Write([]byte(source))
	if err != nil {
		return "", nil
	}
	hashed := h.Sum(nil)
	signatureByte, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signatureByte), nil
}