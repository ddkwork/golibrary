package stream

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max)
}

func RandomNum(min, max int) int {
	rand.NewSource(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// GenXiLa26 生成26个希腊字母
func GenXiLa26() (letters []string) {
	letters = make([]string, 0)
	for i := 'A'; i <= 'Z'; i++ {
		letters = append(letters, string(i))
	}
	return
}
