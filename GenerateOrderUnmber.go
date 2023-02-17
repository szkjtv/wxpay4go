package wxpay

import (
	"math/rand"
	"strconv"
	"time"
)

func getNonceStr() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 10)
}
