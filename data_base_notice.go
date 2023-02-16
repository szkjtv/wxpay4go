package wxpay

import (
	"encoding/json"
	"net/http"
)

//支付回调数据定义
type WeixinPayNotice struct {
	//通知的唯一ID
	ReqId      		string 				`json:"id"`
	//通知创建的时间
	CreateTime      string 				`json:"create_time"`
	//通知的类型，支付成功通知的类型为TRANSACTION.SUCCESS
	Event_type 		string 				`json:"event_type"`
	//通知的资源数据类型，支付成功通知为encrypt-resource
	ResourceType   string 				`json:"resource_type"`
	//回调摘要
	Summary         string 				`json:"summary"`
	//通知资源数据
	Resource 		NoticeResource 		`json:"resource"`
}

type NoticeResource struct {
	//对开启结果数据进行加密的加密算法，目前只支持AEAD_AES_256_GCM
	Algorithm      	string		`json:"algorithm"`
	//数据密文
	Ciphertext     	string		`json:"ciphertext"`
	//附加数据
	AssociatedData 	string		`json:"associated_data"`
	//原始回调类型，为transaction
	original_type  	string		`json:"original_type"`
	//加密使用的随机串
	Nonce          	string		`json:"nonce"`
}

//支付通知返回数据结构
type WxPayNotifyRet struct {
	Return_code 	string 		`json:"code"`
	Return_msg  	string 		`json:"message"`
}

//支付通知返回
func HttpCallBackReturn(w http.ResponseWriter, status int, code string, message string)  {
	var ent WxPayNotifyRet
	ent.Return_code = code
	ent.Return_msg = message
	data, _ := json.Marshal(ent)
	w.WriteHeader(status)
	w.Write(data)
}

