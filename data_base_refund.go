package wxpay

//退款金额信息
type RefundCreateAmount struct {
	//退款金额，单位为分
	Refund 			int 		`json:"refund"`
	//原订单金额，单位为分
	Total 			int 		`json:"total"`
	//退款币种:CNY(人民币)，境内商户号仅支持人民币。
	Currency       	string		`json:"currency"`
}

//微信支付退款申请请求参数
type RefundCreateReq struct {
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no,omitempty"`
	//微信支付订单号
	Transaction_id     	string	`json:"transaction_id,omitempty"`
	//商户系统内部的退款单号
	Out_refund_no     	string	`json:"out_refund_no"`
	//退款原因
	Reason       		string	`json:"reason,omitempty"`
	//退款结果回调的URL，不允许携带查询串
	Notify_url 			string	`json:"notify_url,omitempty"`
	//退款金额信息
	Amount 	RefundCreateAmount 	`json:"amount"`
}

//查询单笔退款参数
type RefundQueryReq struct {
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid,omitempty"`
	//商户系统内部的退款单号
	Out_refund_no     	string	`json:"out_refund_no"`
}

//退款金额信息
type RefundOrderAmount struct {
	//原订单金额，单位为分
	Total 			int 		`json:"total"`
	//退款金额，单位为分
	Refund 			int 		`json:"refund"`
	//用户支付金额，单位为分
	Payer_total 	int 		`json:"payer_total,omitempty"`
	//退款给用户的金额，不包含所有优惠券金额
	Payer_refund 	int 		`json:"payer_refund,omitempty"`
	//应结退款金额
	Settlement_refund 	int 	`json:"settlement_refund"`
	//应结订单金额
	Settlement_total 	int 	`json:"settlement_total"`
	//优惠退款金额
	Discount_refund 	int 	`json:"discount_refund"`
	//退款币种:CNY(人民币)，境内商户号仅支持人民币。
	Currency       	string		`json:"currency"`
}

//微信退款订单信息
type RefundOrderInfo struct {
	//服务商户号(普通商户)
	Mchid          		string	`json:"mchid"`
	//服务商户号（服务商）
	Sp_mchid          	string	`json:"sp_mchid"`
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid"`
	//微信支付退款单号
	Refund_id       	string	`json:"refund_id"`
	//商户系统内部的退款单号
	Out_refund_no     	string	`json:"out_refund_no"`
	//微信支付订单号
	Transaction_id     	string	`json:"transaction_id"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
	//退款渠道 枚举值： - ORIGINAL—原路退款 - BALANCE—退回到余额 - OTHER_BALANCE—原账户异常退到其他余额账户 - OTHER_BANKCARD—原银行卡异常退到其他银行卡 * `ORIGINAL` - 原路退款 * `BALANCE` - 退回到余额 * `OTHER_BALANCE` - 原账户异常退到其他余额账户 * `OTHER_BANKCARD` - 原银行卡异常退到其他银行卡
	Channel 			string `json:"channel"`
	//退款入账账户
	User_received_account    string	`json:"user_received_account"`
	//退款成功时间
	Success_time       	string	`json:"success_time"`
	//退款创建时间
	Create_time       	string	`json:"create_time"`
	//退款状态，UCCESS：退款成功	CLOSED：退款关闭/PROCESSING：退款处理中	ABNORMAL：退款异常
	Status     			string	`json:"status"`
	//退款所使用资金对应的资金账户类型 枚举值：UNSETTLED : 未结算资金AVAILABLE : 可用余额UNAVAILABLE : 不可用余额OPERATION : 运营户BASIC : 基本账户
	Funds_account       string	`json:"funds_account"`
	//订单金额信息
	Amount 	RefundOrderAmount	`json:"amount"`
}

//退款金额信息
type RefundOrderAmountCB struct {
	//退款金额，单位为分
	Refund 			int 		`json:"refund"`
	//原订单金额，单位为分
	Total 			int 		`json:"total"`
	//退款币种:CNY(人民币)，境内商户号仅支持人民币。
	Currency       	string		`json:"currency"`
	//用户支付金额，单位为分
	Payer_total 	int 		`json:"payer_total,omitempty"`
	//退款给用户的金额，不包含所有优惠券金额
	Payer_refund 	int 		`json:"payer_refund,omitempty"`
}

//微信退款订单信息
type RefundOrderInfoCB struct {
	//服务商户号(普通商户)
	Mchid          		string	`json:"mchid"`
	//服务商户号（服务商）
	Sp_mchid          	string	`json:"sp_mchid"`
	//子商户的商户号（服务商）
	Sub_mchid       	string	`json:"sub_mchid"`
	//微信支付订单号
	Transaction_id     	string	`json:"transaction_id"`
	//商户系统内部订单号
	Out_trade_no     	string	`json:"out_trade_no"`
	//微信支付退款单号
	Refund_id       	string	`json:"refund_id"`
	//商户系统内部的退款单号
	Out_refund_no     	string	`json:"out_refund_no"`
	//退款状态，SUCCESS：退款成功 CLOSE：退款关闭 ABNORMAL：退款异常
	Refund_status     	string	`json:"refund_status"`
	//退款成功时间
	Success_time       	string	`json:"success_time"`
	//退款入账账户
	User_received_account	string	`json:"user_received_account"`
	//订单金额信息
	Amount 		RefundOrderAmountCB `json:"amount"`
}

