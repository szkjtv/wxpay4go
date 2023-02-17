package wxpay

import (
	"math/rand"
	"strconv"
	"time"
)

// 根据时间戳生成订单号
func GetNonceStr() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 10)
}
