package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MchParam struct {
	//商户对应的appid
	Appid string
	//商户号
	Mchid string
	//商户的API v3密钥
	MchAPIKey string
	//商户API私钥
	MchPrivateKey *rsa.PrivateKey
	//商户 API 证书
	MchCertificate *x509.Certificate
}

// 支付回调请求处理，返回订单信息
func (ent *MchParam) ParsePayCallBack(r *http.Request) (WxPayInfo, error) {
	var pay_info WxPayInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return pay_info, err
	}
	if len(body) < 1 {
		err := fmt.Errorf("读取http body失败")
		return pay_info, err
	}

	//读取签名验证所需的参数
	var sing_param WxSignParam
	err = sing_param.GetFromCallBackRequest(r, string(body))
	if err != nil {
		return pay_info, err
	}

	//获取平台证书，并进行签名验证
	plat_certificate := GetPlatCertificate(ent, sing_param.CertSerial)
	err = ResponseValidate(&sing_param, plat_certificate)
	if err != nil {
		return pay_info, err
	}

	//body数据解析
	var ent_cb WeixinPayNotice
	if err = json.Unmarshal(body, &ent_cb); err != nil {
		return pay_info, err
	}

	//数据解密
	decryptBytes, err := DecryptAES256GCM(
		ent.MchAPIKey,
		ent_cb.Resource.AssociatedData,
		ent_cb.Resource.Nonce,
		ent_cb.Resource.Ciphertext)
	if err != nil {
		return pay_info, err
	}

	//支付订单数据解析
	err = json.Unmarshal([]byte(decryptBytes), &pay_info)
	if err != nil {
		return pay_info, err
	}
	return pay_info, nil
}

// 退款回调请求处理，返回退款信息
func (ent *MchParam) ParseRefundCallBack(r *http.Request) (RefundOrderInfoCB, error) {
	var ret_info RefundOrderInfoCB
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ret_info, err
	}
	if len(body) < 1 {
		err := fmt.Errorf("读取http body失败")
		return ret_info, err
	}

	//读取签名验证所需的参数
	var sing_param WxSignParam
	err = sing_param.GetFromCallBackRequest(r, string(body))
	if err != nil {
		return ret_info, err
	}

	//获取平台证书，并进行签名验证
	plat_certificate := GetPlatCertificate(ent, sing_param.CertSerial)
	err = ResponseValidate(&sing_param, plat_certificate)
	if err != nil {
		return ret_info, err
	}

	//body数据解析
	var ent_cb WeixinPayNotice
	if err = json.Unmarshal(body, &ent_cb); err != nil {
		return ret_info, err
	}

	//数据解密
	decryptBytes, err := DecryptAES256GCM(
		ent.MchAPIKey,
		ent_cb.Resource.AssociatedData,
		ent_cb.Resource.Nonce,
		ent_cb.Resource.Ciphertext)
	if err != nil {
		return ret_info, err
	}

	//支付订单数据解析
	err = json.Unmarshal([]byte(decryptBytes), &ret_info)
	if err != nil {
		return ret_info, err
	}
	return ret_info, nil
}

// 支付通知返回
func (ent *MchParam) HttpCallBackReturn(w http.ResponseWriter, status int, code string, message string) {
	var ret WxPayNotifyRet
	ret.Return_code = code
	ret.Return_msg = message
	data, _ := json.Marshal(ret)
	w.WriteHeader(status)
	w.Write(data)
}

type WxPayer interface {
	// App支付统一下单
	CreateOrderApp(data OsAppOrderData) (RawAppPayParam, error)
	// JsApi支付统一下单
	CreateOrderJsApi(data JsApiOrderData) (WxAppPayParam, error)
	// Native支付统一下单
	CreateOrderNative(data NativeOrderData) (string, error)
	// H5支付统一下单
	CreateOrderH5(data WxH5OrderData) (string, error)

	// 支付订单查询
	QueryByTransactionId(transaction_id string, mchid string) (WxPayInfo, error)
	// 支付订单查询
	QueryByOutTradeNo(out_trade_no string, mchid string) (WxPayInfo, error)
	// 关闭支付订单
	CloseOrder(out_trade_no string, mchid string) error
	// 提交退款申请
	RefundOrder(data RefundCreateReq) (RefundOrderInfo, error)
	// 查询单笔退款
	QueryRefundOrder(out_refund_no string, mchid string) (RefundOrderInfo, error)
}

/*
支付下单
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/jsapi（小程序）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/jsapi（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/app（APP）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/h5（H5）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/native（Native）
订单号查询
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/{transaction_id}（小程序）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/{transaction_id}（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/{transaction_id}（APP）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/{transaction_id}（H5）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/{transaction_id}（Native）
关闭订单
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/{out_trade_no}/close（小程序）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/{out_trade_no}/close（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/{out_trade_no}/close（APP）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/{out_trade_no}/close（H5）
https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/{out_trade_no}/close（Native）
申请退款
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（小程序）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（JSAPI）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（APP）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（H5）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（Native）
查询单笔退款
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（小程序）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（JSAPI）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（APP）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（H5）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（Native）

支付下单
https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi（小程序）
https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/transactions/app（APP）
https://api.mch.weixin.qq.com/v3/pay/transactions/h5（H5）
https://api.mch.weixin.qq.com/v3/pay/transactions/native（Native）
商户号查询
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}（小程序）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}（APP）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}（H5）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}（Native）
关闭订单
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}/close（小程序）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}/close（JSAPI）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}/close（APP）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}/close（H5）
https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/{out_trade_no}/close（Native）
申请退款
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（小程序）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（JSAPI）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（APP）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（H5）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds（Native）
查询单笔退款
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（小程序）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（JSAPI）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（APP）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（H5）
https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/{out_refund_no}（Native）
*/
