package mathutil

import (
	"math/rand"
	"time"
)

func GenRandomInt(min, max int) int {
	if max-min <= 0 {
		return 0
	}
	return rand.New(rand.NewSource(time.Now().Unix())).Intn(max-min) + min
}

func GenRandomUint32() uint32 {
	return rand.New(rand.NewSource(time.Now().Unix())).Uint32()
}
