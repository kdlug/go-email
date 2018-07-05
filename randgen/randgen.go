package randgen

import (
	"fmt"
	"math/rand"
	"time"
)

var chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randChars(n int) string {

	b := make([]byte, n)
	for i := range b {
		var seed = rand.NewSource(time.Now().UnixNano())
		r := rand.New(seed)
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}

func wrap(s string, pre string, suf string) string {
	res := ""
	res += fmt.Sprintf("%s%s%s", pre, s, suf)

	return res

}

// GenerateStr generates random string and wraps it using prefix and sufix
func GenerateStr(n int, prefix string, suffix string) string {
	s := randChars(n)
	s = wrap(s, prefix, suffix)

	return s
}
