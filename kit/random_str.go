package kit

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

// https://xie.infoq.cn/article/f274571178f1bbe6ff8d974f3

type StrRandomElement uint8

func (s StrRandomElement) IsSet(s1 StrRandomElement) bool {
	return s&s1 == s1
}

const (
	ALLElement = Number | LowerLetter | UpperLetter | SpecialSymbol

	Number StrRandomElement = 1 << iota
	LowerLetter
	UpperLetter
	SpecialSymbol
)

// random source
const (
	lowerLetters   = "abcdefghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers        = "0123456789"
	specialSymbols = `~!@#$%^&*()_-+={[}]|\:;"'<,>.?/`
)

type RandomStr struct {
	randSrc rand.Source
	strSrc  string
}

func NewRandomStr(elements StrRandomElement) *RandomStr {
	rs := &RandomStr{
		randSrc: rand.NewSource(time.Now().Unix()),
		strSrc:  buildRandomStrSource(elements),
	}

	return rs
}

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func (r *RandomStr) Generate(length int) string {
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := length-1, r.randSrc.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.randSrc.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(r.strSrc) {
			b[i] = r.strSrc[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}

	// r.randSrc = rand.NewSource(time.Now().Unix())
	return *(*string)(unsafe.Pointer(&b))
}

func buildRandomStrSource(elements StrRandomElement) string {
	var strBuilder strings.Builder
	if elements.IsSet(Number) {
		strBuilder.WriteString(numbers)
	}
	if elements.IsSet(LowerLetter) {
		strBuilder.WriteString(lowerLetters)
	}
	if elements.IsSet(UpperLetter) {
		strBuilder.WriteString(upperLetters)
	}
	if elements.IsSet(SpecialSymbol) {
		strBuilder.WriteString(specialSymbols)
	}

	return strBuilder.String()
}
