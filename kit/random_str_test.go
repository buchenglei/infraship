package kit

import (
	"testing"
	"time"
)

func TestBuildRandomStrSource(t *testing.T) {
	eles := UpperLetter | SpecialSymbol
	t.Log(buildRandomStrSource(eles))
}

func TestRandomStrGenerate(t *testing.T) {
	eles := Number
	rs := NewRandomStr(eles)

	for i := 20; i >= 0; i-- {
		time.Sleep(10 * time.Millisecond)
		t.Log(rs.Generate(6))
	}
}
