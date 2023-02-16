package wxpay

import (
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//根据证书序列号获取平台证书
var map_plat_cert = make(map[string]*x509.Certificate)
var lock_plat_cert sync.RWMutex
func GetPlatCertificate(ctx *MchParam, serial_no string) *x509.Certificate {
	lock_plat_cert.Lock()
	defer lock_plat_cert.Unlock()

	cert, ok := map_plat_cert[serial_no]
	if !ok {
		//调用证书获取接口
		downloadPlatCertificate(ctx)
		cert = map_plat_cert[serial_no]
		if cert == nil {
			log.Println("GetPlatCertificate error !!!")
		}
	}
	return cert
}

//微信支付平台证书接口返回结果定义
type WxCertRet struct {
	Data 	[]WxCertRetItem 	`json:"data"`
}

type WxCertRetItem struct {
	EffectiveTime string        `json:"effective_time"`
	Certificate   WxCertificate `json:"encrypt_certificate"`
	ExpireTime    string        `json:"expire_time"`
	SerialNo      string        `json:"serial_no"`
}

type WxCertificate struct {
	Algorithm      string		`json:"algorithm"`
	AssociatedData string		`json:"associated_data"`
	Ciphertext     string		`json:"ciphertext"`
	Nonce          string		`json:"nonce"`
}

//获取微信支付平台证书
func downloadPlatCertificate(ctx *MchParam) error {
	const publicKeyUrl = "https://api.mch.weixin.qq.com/v3/certificates"
	log.Println(publicKeyUrl)
	token, err := CreateAuthorization(ctx, http.MethodGet, publicKeyUrl, "")
	if err != nil {
		log.Println(err)
		return err
	}

	request, err := http.NewRequest(http.MethodGet, publicKeyUrl, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	request.Header.Add("Authorization", token)
	request.Header.Add("User-Agent", "go pay sdk")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	var tokenResponse WxCertRet
	if err = json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		log.Println(err)
		return err
	}

	//使用AEAD_AES_256_GCM算法进行解密
	for _, encrypt_cert := range tokenResponse.Data {
		decryptBytes, err := DecryptAES256GCM(
			ctx.MchAPIKey,
			encrypt_cert.Certificate.AssociatedData,
			encrypt_cert.Certificate.Nonce,
			encrypt_cert.Certificate.Ciphertext)
		if err != nil {
			log.Println(err)
			return err
		}

		cert_ret, _ := LoadCertificate(decryptBytes)
		if cert_ret != nil {
			serial_no := encrypt_cert.SerialNo
			log.Printf("certificate:%s\n", serial_no)
			map_plat_cert[serial_no] = cert_ret
		}
	}
	return nil
}

