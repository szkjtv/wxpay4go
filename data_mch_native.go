package wxpay

//统一下单请求参数
type NativeOrderData struct {
	//子商户用ID（服务商）
	Sub_appid       	string	`json:"sub_appid,omitempty"`
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//商品描述
	Description       	string	`json:"description"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
	//交易结束时间
	Time_expire       	string	`json:"time_expire,omitempty"`
	//附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用
	Attach       	   	string	`json:"attach,omitempty"`
	//通知URL必须为直接可访问的URL，不允许携带查询串
	Notify_url 			string	`json:"notify_url"`
	//订单优惠标记
	Goods_tag       	string	`json:"goods_tag,omitempty"`
	//订单金额信息
	Amount OrderCreateAmount 		`json:"amount"`
	//统一下单结算信息
	Settle_info OrderCreateSettle 	`json:"settle_info"`
}

//统一下单请求参数
type NativeOrderCreateReq struct {
	//应用ID(普通商户)
	Appid       		string	`json:"appid,omitempty"`
	//直连商户号(普通商户)
	Mchid          		string	`json:"mchid,omitempty"`
	//服务商用ID（服务商）
	Sp_appid       		string	`json:"sp_appid,omitempty"`
	//服务商户号（服务商）
	Sp_mchid          	string	`json:"sp_mchid,omitempty"`
	//订单数据
	NativeOrderData
}

//统一下单返回参数
type NativeOrderCreateRet struct {
	//详细错误码
	Return_code 		string 	`json:"code"`
	//错误描述
	Return_msg  		string 	`json:"message"`
	//二维码链接,此URL用于生成支付二维码，然后提供给用户扫码支付。
	Code_url   			string 	`json:"code_url"`
}
