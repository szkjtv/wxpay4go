package wxpay

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOrderUnmber() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 10)
}
