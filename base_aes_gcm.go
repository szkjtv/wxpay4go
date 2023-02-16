package wxpay

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//使用 AEAD_AES_256_GCM 算法进行解密
//你可以使用此算法完成微信支付平台证书和回调报文解密
func DecryptAES256GCM(aesKey, associatedData, nonce, ciphertext string) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	dataBytes, err := gcm.Open(nil, []byte(nonce), decodedCiphertext, []byte(associatedData))
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
