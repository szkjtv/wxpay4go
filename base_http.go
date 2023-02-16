package wxpay

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 生成请求报文头中的Authorization信息
// method:HTTP请求的方法（GET,POST）
// rawUrl:请求的绝对URL，用于获取除域名部分得到参与签名的URL
// signBody:请求报文主体
func CreateAuthorization(ctx *MchParam, method string, rawUrl string, signBody string) (authorization string, err error) {
	timestamp := time.Now().Unix()
	url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	nonce, err := GenerateNonce()
	if err != nil {
		return "", err
	}

	SignatureMessageFormat := "%s\n%s\n%d\n%s\n%s\n"
	message := fmt.Sprintf(SignatureMessageFormat, method, url.RequestURI(), timestamp, nonce, signBody)
	//gowf.LogD(message)
	signatureResult, err := SignSHA256WithRSA(ctx.MchPrivateKey, message)
	if err != nil {
		return "", err
	}

	certSerialNo := fmt.Sprintf("%X", ctx.MchCertificate.SerialNumber)
	HeaderAuthorizationFormat := "WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\""
	authorization = fmt.Sprintf(HeaderAuthorizationFormat, ctx.Mchid, nonce, timestamp, certSerialNo, signatureResult)
	return authorization, nil
}

type WxSignParam struct {
	Timestamp  string //微信支付回包时间戳
	Nonce      string //微信支付回包随机字符串
	Signature  string //微信支付回包签名信息
	CertSerial string //微信支付回包平台序列号
	RequestId  string //微信支付回包请求ID
	Body       string
}

// 构造验签名串
func (ent *WxSignParam) BuildResponseMessage() string {
	message := fmt.Sprintf("%s\n%s\n%s\n",
		ent.Timestamp, ent.Nonce, ent.Body)
	return message
}

// 对微信支付应答报文进行验证
// resp_ent:微信支付应答报文数据
// certificate:微信支付平台证书中，使用微信支付平台证书中的公钥验签
func ResponseValidate(sign_param *WxSignParam, certificate *x509.Certificate) error {
	message := sign_param.BuildResponseMessage()
	signature, err := base64.StdEncoding.DecodeString(sign_param.Signature)
	if err != nil {
		return fmt.Errorf("base64 decode string wechat pay signature err:%s", err.Error())
	}

	hashed := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(certificate.PublicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], []byte(signature))
	if err != nil {
		return fmt.Errorf("verifty signature with public key err:%s", err.Error())
	}

	return nil
}

func (ent *WxSignParam) GetFromHttpResponse(response *http.Response, body string) error {
	ent.RequestId = strings.TrimSpace(response.Header.Get("Request-Id"))
	if ent.RequestId == "" {
		return fmt.Errorf("empty Request-Id")
	}
	ent.CertSerial = strings.TrimSpace(response.Header.Get("Wechatpay-Serial"))
	if ent.CertSerial == "" {
		return fmt.Errorf("empty WechatPaySerial")
	}
	ent.Signature = strings.TrimSpace(response.Header.Get("Wechatpay-Signature"))
	if ent.Signature == "" {
		return fmt.Errorf("empty WechatPaySignature")
	}
	ent.Timestamp = strings.TrimSpace(response.Header.Get("Wechatpay-Timestamp"))
	if ent.Timestamp == "" {
		return fmt.Errorf("empty WechatPayTimestamp")
	}
	ent.Nonce = strings.TrimSpace(response.Header.Get("Wechatpay-Nonce"))
	if ent.Nonce == "" {
		return fmt.Errorf("empty Wechatpay-Nonce")
	}
	ent.Body = body
	if len(body) < 1 {
		return fmt.Errorf("read response body")
	}
	return nil
}

func (ent *WxSignParam) GetFromCallBackRequest(r *http.Request, body string) error {
	ent.CertSerial = strings.TrimSpace(r.Header.Get("Wechatpay-Serial"))
	if ent.CertSerial == "" {
		return fmt.Errorf("empty WechatPaySerial")
	}
	ent.Signature = strings.TrimSpace(r.Header.Get("Wechatpay-Signature"))
	if ent.Signature == "" {
		return fmt.Errorf("empty WechatPaySignature")
	}
	ent.Timestamp = strings.TrimSpace(r.Header.Get("Wechatpay-Timestamp"))
	if ent.Timestamp == "" {
		return fmt.Errorf("empty WechatPayTimestamp")
	}
	ent.Nonce = strings.TrimSpace(r.Header.Get("Wechatpay-Nonce"))
	if ent.Nonce == "" {
		return fmt.Errorf("empty Wechatpay-Nonce")
	}
	ent.Body = body
	if len(body) < 1 {
		return fmt.Errorf("read response body")
	}
	return nil
}

// 通过http GET方法调用支付接口
// ctx:商户上下文对象
// url:接口的调用地址
func WxPayGetV3(ctx *MchParam, url string) (string, error) {
	log.Println(url)
	if ctx.MchCertificate == nil {
		return "", fmt.Errorf("没有证书文件")
	}
	if ctx.MchPrivateKey == nil {
		return "", fmt.Errorf("没有密钥文件")
	}
	token, err := CreateAuthorization(ctx, http.MethodGet, url, "")
	if err != nil {
		log.Println(err)
		return "", err
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("User-Agent", "go pay sdk")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		err := fmt.Errorf("status:%d;msg=%s", resp.StatusCode, string(result))
		log.Println(err)
		return string(result), err
	}

	var sign_param WxSignParam
	err = sign_param.GetFromHttpResponse(resp, string(result))
	if err != nil {
		log.Println(err)
		return string(result), err
	}

	//get plat certificate
	plat_certificate := GetPlatCertificate(ctx, sign_param.CertSerial)
	if plat_certificate == nil {
		err := fmt.Errorf("plat_certificate get error:%s", sign_param.CertSerial)
		log.Println(err)
		return string(result), err
	}

	//Validate WechatPay Signature
	err = ResponseValidate(&sign_param, plat_certificate)
	if err != nil {
		log.Println(err)
		return string(result), err
	}

	log.Println(string(result))
	return string(result), nil
}

// 通过http POST方法调用支付接口
// ctx:商户上下文对象
// url:接口的调用地址
// data:报文主体的JSON数据
func WxPayPostV3(ctx *MchParam, url string, data []byte) (string, error) {
	log.Println(url)
	log.Println(string(data))

	if ctx.MchCertificate == nil {
		return "", fmt.Errorf("没有证书文件")
	}
	if ctx.MchPrivateKey == nil {
		return "", fmt.Errorf("没有密钥文件")
	}
	token, err := CreateAuthorization(ctx, http.MethodPost, url, string(data))
	if err != nil {
		log.Println(err)
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("User-Agent", "go pay sdk")
	request.Header.Add("Content-type", "application/json;charset='utf-8'")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		var pret WxPayNotifyRet
		err = json.Unmarshal([]byte(result), &pret)
		if err != nil {
			err = fmt.Errorf("status:%d;msg=%s", resp.StatusCode, string(result))
		} else {
			err = errors.New(pret.Return_msg)
		}
		log.Println(err)
		return string(result), err
	}

	var sign_param WxSignParam
	err = sign_param.GetFromHttpResponse(resp, string(result))
	if err != nil {
		log.Println(err)
		return string(result), err
	}

	//get plat certificate
	plat_certificate := GetPlatCertificate(ctx, sign_param.CertSerial)
	if plat_certificate == nil {
		err := fmt.Errorf("plat_certificate get error:%s", sign_param.CertSerial)
		log.Println(err)
		return string(result), err
	}

	//Validate WechatPay Signature
	err = ResponseValidate(&sign_param, plat_certificate)
	if err != nil {
		log.Println(err)
		return string(result), err
	}

	log.Println(string(result))
	return string(result), nil
}
