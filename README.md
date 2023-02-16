### 介绍
微信支付SDK，基于全新的微信支付APIv3来实现。支持小程序支付、JSAPI支付、Native支付、APP支付、H5支付，支持直连商户模式和服务商模式。

### 安装说明
go get gitee.com/haming123/wxpay4go

### 快速上手
1. 商户对象初始化
```go
package main

import (
	"log"
	"wxpay"
)

var MchCtx *wxpay.MchWxApp = nil
func MchCtxInit() error {
	pem_cert, err := wxpay.LoadCertificateWithPath("/path/to/merchant/apiclient_cert.pem")
	if err != nil {
		return err
	}
	pem_key, err := wxpay.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		return err
	}
	MchCtx = wxpay.NewMchWxApp("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
	return nil
}

func main() {
	//......
	err := MchCtxInit()
	if err != nil {
		log.Print(err)
		return
	}
	//......
}
```

2. 支付下单
```go
func HandlerPayCreateDemo(w http.ResponseWriter, r *http.Request) {
	//................(其他业务逻辑)
	var data_req wxpay.WxAppOrderData
	data_req.Description = "业务描述"
	data_req.Out_trade_no = "业务订单号"
	data_req.Notify_url = "接收支付通知的URL地址"
	data_req.Sub_mchid = "子商户号"
	data_req.Sub_appid = "小程序的APPID"
	data_req.Payer.Sub_openid = "用户在小程序中的OPENID"
	data_req.Amount.Total = 20
	param, err := MchCtx.CreateOrder(data_req)
	if err!= nil {
		log.Print(err)
		return
	}
	//................(其他业务逻辑)
}
```

3. 支付订单查询
```go
func HandlerPayQueryDemo(w http.ResponseWriter, r *http.Request) {
	//................(其他业务逻辑)
	out_trade_no := "业务订单号"
	sub_mchid := "子商户号"
	result, err := MchCtx.QueryByOutTradeNo(out_trade_no, sub_mchid)
	if err!= nil {
		log.Print(err)
		return
	}
	//................(其他业务逻辑)
}
```

4. 支付回调处理
```go
func HandlerPayCallBack(w http.ResponseWriter, r *http.Request) {
	//解析回调数据
	req_info, err := MchCtx.ParsePayCallBack(r)
	if err != nil {
		wxpay.HttpCallBackReturn(w, 500,"FAIL", "FAIL")
		return
	}

	//................(其他业务逻辑)

	//回调返回
	wxpay.HttpCallBackReturn(w, 200, "SUCCESS", "SUCCESS")
}
```

5. 退款申请
```go
func HandlerRefundDemo(w http.ResponseWriter, r *http.Request) {
	//................(其他业务逻辑)
	var data wxpay.RefundCreateReq
	data.Sub_mchid = "子商户号"
	data.Out_trade_no = "业务订单号"
	data.Out_refund_no = "商户的退款单号"
	data.Amount.Total = 20
	data.Amount.Refund = 20
	data.Notify_url = "退款结果回调的URL"
	result, err := MchCtx.RefundOrder(data)
	if err!= nil {
		log.Print(err)
	}
	//................(其他业务逻辑)
}
```

6. 退款回调处理
```go
func HandlerRefundCallBack(w http.ResponseWriter, r *http.Request) {
	//解析回调数据
	req_info, err := MchCtx.ParseRefundCallBack(r)
	if err != nil {
		wxpay.HttpCallBackReturn(w, 500,"FAIL", "FAIL")
		return
	}

	//................(其他业务逻辑)

	//回调返回
	wxpay.HttpCallBackReturn(w, 200, "SUCCESS", "SUCCESS")
}
```

### 其他支付产品的使用
1. 直连商户商户对象初始化
```go
//小程序
MchCtx = wxpay.NewMchWxApp("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//JsAPi
MchCtx = wxpay.NewMchJsApi("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//H5
MchCtx = wxpay.NewMchH5("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//Native
MchCtx = wxpay.NewMchNative("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//App
MchCtx = wxpay.NewMchApp("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
```

2. 服务商商户对象初始化
```go
//小程序
MchCtx = wxpay.NewMidWxApp("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//JsAPi
MchCtx = wxpay.NewMidJsApi("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//H5
MchCtx = wxpay.NewMidH5("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//Native
MchCtx = wxpay.NewMidNative("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
//App
MchCtx = wxpay.NewMidApp("mchID_string", "appID_string", "mchAPIv3Key_string", pem_cert, pem_key)
```
