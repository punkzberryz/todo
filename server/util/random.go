package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvxyz"

func getRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomBool() bool {
	return getRand().Float32() > 0.5
}

func RandomInt(min, max int64) int64 {
	//example for 6-digit RandomInt(100000,999999)
	return min + getRand().Int63n(max-min+1) //0->max-min
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[getRand().Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
