package wxpay

//统一下单结算信息
type OrderCreateSettle struct {
	//是否指定分账，枚举值:true(是), false(否)
	Profit_sharing 	bool 		`json:"profit_sharing,omitempty"`
}

//统一下单订单金额信息
type OrderCreateAmount struct {
	//订单总金额，单位为分。
	Total 			int 		`json:"total"`
	//货币类型:CNY(人民币)，境内商户号仅支持人民币。
	Currency       	string		`json:"currency,omitempty"`
}

//支付订单查询参数
type QueryByTransactionIdReq struct {
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//微信支付订单号
	Transaction_id     	string	`json:"transaction_id"`
}

//支付订单查询参数
type QueryByOutTradeNoReq struct {
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
}

//关闭支付订单参数
type CloseOrderReq struct {
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
}

//支付者信息
type OrderPayer struct {
	//用用户在直连商户appid下的唯一标识
	Openid       		string		`json:"openid,omitempty"`
	//用户在服务商appid下的唯一标识
	Sp_openid       	string		`json:"sp_openid,omitempty"`
	//用户在子商户appid下的唯一标识。若传sub_openid，那sub_appid必填
	Sub_openid       	string		`json:"sub_openid,omitempty"`
}

//订单金额信息
type OrderAmount struct {
	//订单总金额，单位为分
	Total 			int 		`json:"total"`
	//用户支付金额，单位为分
	Payer_total 	int 		`json:"payer_total"`
	//货币类型:CNY(人民币)，境内商户号仅支持人民币。
	Currency       	string		`json:"currency"`
	//用户支付币种
	Payer_currency  string		`json:"payer_currency"`
}

//微信支付订单在信息
type WxPayInfo struct {
	//应用ID(普通商户)
	Appid       		string	`json:"appid"`
	//服务商户号(普通商户)
	Mchid          		string	`json:"mchid"`
	//服务商用ID（服务商）
	Sp_appid       		string	`json:"sp_appid"`
	//服务商户号（服务商）
	Sp_mchid          	string	`json:"sp_mchid"`
	//子商户用ID（服务商）
	Sub_appid       	string	`json:"sub_appid"`
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
	//微信支付订单号
	Transaction_id     	string	`json:"transaction_id"`
	//交易类型，枚举值：JSAPI
	Trade_type     		string	`json:"trade_type"`
	//交易状态，SUCCESS：支付成功 REFUND：转入退款 NOTPAY：未支付 CLOSED：已关闭 PAYERROR：支付失败
	Trade_state     	string	`json:"trade_state"`
	//交易状态描述
	Trade_state_desc    string	`json:"trade_state_desc"`
	//付款银行
	Bank_type     		string	`json:"bank_type"`
	//附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用
	Attach       	   	string	`json:"attach"`
	//支付完成时间
	Success_time       	string	`json:"success_time"`
	//订单金额信息
	Amount 		OrderAmount 	`json:"amount"`
	//支付者信息
	Payer 		OrderPayer 		`json:"payer"`
}
