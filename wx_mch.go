package wxpay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type WxPayMch struct {
	MchParam
}

func NewWxPayMch(mchid, appid, mch_api_key string, pem_cert *x509.Certificate, pem_key *rsa.PrivateKey) *WxPayMch {
	ent := &WxPayMch{}
	ent.Appid = appid
	ent.Mchid = mchid
	ent.MchAPIKey = mch_api_key
	ent.MchPrivateKey = pem_key
	ent.MchCertificate = pem_cert
	return ent
}

// 支付订单查询
// transaction_id:微信支付订单号
// mchid:商户号
func (ent *WxPayMch) QueryByTransactionId(transaction_id string, mchid string) (WxPayInfo, error) {
	var par_info WxPayInfo
	url := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/pay/transactions/id/%s?mchid=%s", transaction_id, mchid)
	result, err := WxPayGetV3(&ent.MchParam, url)
	if err != nil {
		fmt.Println(err)
		return par_info, err
	}
	err = json.Unmarshal([]byte(result), &par_info)
	if err != nil {
		fmt.Println(err)
		return par_info, err
	}

	return par_info, nil
}

// 支付订单查询
// out_trade_no:业务订单号
// mchid:商户号
func (ent *WxPayMch) QueryByOutTradeNo(out_trade_no string, mchid string) (WxPayInfo, error) {
	var par_info WxPayInfo
	url := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/%s?mchid=%s", out_trade_no, mchid)
	result, err := WxPayGetV3(&ent.MchParam, url)
	if err != nil {
		return par_info, err
	}
	err = json.Unmarshal([]byte(result), &par_info)
	if err != nil {
		return par_info, err
	}

	return par_info, nil
}

// 关闭支付订单
// out_trade_no:业务订单号
// mchid:商户号
func (ent *WxPayMch) CloseOrder(out_trade_no string, mchid string) error {
	type CloseOrderReq struct {
		Mchid string `json:"mchid,omitempty"`
	}
	var preq CloseOrderReq
	preq.Mchid = mchid
	data_body, _ := json.Marshal(preq)

	url := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/%s/close", out_trade_no)
	_, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return err
	}
	return nil
}

// 提交退款申请
func (ent *WxPayMch) RefundOrder(data RefundCreateReq) (RefundOrderInfo, error) {
	data.Sub_mchid = ""
	if data.Amount.Currency == "" {
		data.Amount.Currency = "CNY"
	}
	data_body, _ := json.Marshal(data)

	var pret RefundOrderInfo
	const url = "https://api.mch.weixin.qq.com/v3/refund/domestic/refunds"
	result, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return pret, err
	}
	err = json.Unmarshal([]byte(result), &pret)
	if err != nil {
		return pret, err
	}
	return pret, nil
}

// 查询单笔退款
// out_refund_no:商户系统内部的退款单号
// mchid:商户号
func (ent *WxPayMch) QueryRefundOrder(out_refund_no string, mchid string) (RefundOrderInfo, error) {
	var ret_info RefundOrderInfo
	url := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/%s", out_refund_no)
	result, err := WxPayGetV3(&ent.MchParam, url)
	if err != nil {
		fmt.Println(err)
		return ret_info, err
	}
	err = json.Unmarshal([]byte(result), &ret_info)
	if err != nil {
		fmt.Println(err)
		return ret_info, err
	}

	return ret_info, nil
}

// JsApi支付统一下单
// data:支付订单信息
// sub_mch_id:子商户号, 传空格即可
func (ent *WxPayMch) CreateOrderJsApi(data JsApiOrderData) (WxAppPayParam, error) {
	var preq JsApiOrderCreateReq
	preq.Appid = ent.Appid
	preq.Mchid = ent.Mchid
	preq.JsApiOrderData = data
	preq.JsApiOrderData.Sub_mchid = ""
	preq.JsApiOrderData.Sub_appid = ""
	data_body, _ := json.Marshal(preq)
	logPrint(BytesToString(data_body))

	var param_ent WxAppPayParam
	var pret JsApiOrderCreateRet
	const url = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
	result, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return param_ent, err
	}

	err = json.Unmarshal([]byte(result), &pret)
	if err != nil {
		return param_ent, err
	}

	if pret.Prepay_id == "" {
		return param_ent, errors.New(pret.Return_msg)
	}

	//小程序客户端支付参数
	param_ent.Appid = preq.Appid
	param_ent.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	param_ent.NonceStr, _ = GenerateNonce()
	param_ent.Prepay_id = pret.Prepay_id
	param_ent.Package = "prepay_id=" + pret.Prepay_id
	param_ent.SignType = "RSA"
	param_ent.PaySign, _ = param_ent.GenPaySignV3(ent.MchPrivateKey)
	return param_ent, nil
}

// Native支付统一下单
func (ent *WxPayMch) CreateOrderNative(data NativeOrderData) (string, error) {
	var preq NativeOrderCreateReq
	preq.Appid = ent.Appid
	preq.Mchid = ent.Mchid
	preq.NativeOrderData = data
	preq.Sp_appid = ""
	preq.Sp_mchid = ""
	preq.NativeOrderData.Sub_mchid = ""
	preq.NativeOrderData.Sub_appid = ""
	data_body, _ := json.Marshal(preq)
	logPrint(BytesToString(data_body))

	var pret NativeOrderCreateRet
	const url = "https://api.mch.weixin.qq.com/v3/pay/transactions/native"
	result, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(result), &pret)
	if err != nil {
		return "", err
	}
	if pret.Code_url == "" {
		return "", errors.New(pret.Return_msg)
	}

	return pret.Code_url, nil
}

// H5支付统一下单
func (ent *WxPayMch) CreateOrderH5(data WxH5OrderData) (string, error) {
	var preq WxH5OrderCreateReq
	preq.Appid = ent.Appid
	preq.Mchid = ent.Mchid
	preq.WxH5OrderData = data
	preq.Sp_appid = ""
	preq.Sp_mchid = ""
	preq.WxH5OrderData.Sub_mchid = ""
	preq.WxH5OrderData.Sub_appid = ""
	if preq.Scene_info.H5_info.SceneType == "" {
		preq.Scene_info.H5_info.SceneType = "Wap"
	}
	data_body, _ := json.Marshal(preq)
	logPrint(BytesToString(data_body))

	var pret WxH5OrderCreateRet
	const url = "https://api.mch.weixin.qq.com/v3/pay/transactions/h5"
	result, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(result), &pret)
	if err != nil {
		return "", err
	}
	if pret.H5_url == "" {
		return "", errors.New(pret.Return_msg)
	}

	return pret.H5_url, nil
}

// App支付统一下单
func (ent *WxPayMch) CreateOrderApp(data OsAppOrderData) (RawAppPayParam, error) {
	var preq OsAppOrderCreateReq
	preq.Appid = ent.Appid
	preq.Mchid = ent.Mchid
	preq.OsAppOrderData = data
	preq.Sp_appid = ""
	preq.Sp_mchid = ""
	preq.OsAppOrderData.Sub_mchid = ""
	preq.OsAppOrderData.Sub_appid = ""
	data_body, _ := json.Marshal(preq)
	logPrint(BytesToString(data_body))

	var param_ent RawAppPayParam
	var pret WxAppOrderCreateRet
	const url = "https://api.mch.weixin.qq.com/v3/pay/transactions/app"
	result, err := WxPayPostV3(&ent.MchParam, url, data_body)
	if err != nil {
		return param_ent, err
	}

	err = json.Unmarshal([]byte(result), &pret)
	if err != nil {
		return param_ent, err
	}

	if pret.Prepay_id == "" {
		return param_ent, errors.New(pret.Return_msg)
	}

	//小程序客户端支付参数
	param_ent.Appid = preq.Appid
	param_ent.PartnerId = ent.Mchid
	param_ent.Prepay_id = pret.Prepay_id
	param_ent.Package = "Sign=WXPay"
	param_ent.NonceStr, _ = GenerateNonce()
	param_ent.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	param_ent.PaySign, _ = param_ent.GenPaySignV3(ent.MchPrivateKey)
	return param_ent, nil
}
