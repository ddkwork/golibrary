package random

import (
	"math/rand"
	"time"
)

type (
	Interface interface {
		RandomNum(min, max int) int //生成指定长度的随机数
	}
	object struct{}
)

func New() Interface { return &object{} }

func (o *object) RandomNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
