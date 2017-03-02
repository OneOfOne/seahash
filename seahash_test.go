package seahash

import (
	"bytes"
	"testing"

	"github.com/OneOfOne/xxhash"
)

var testBuff = bytes.Repeat([]byte("to be or not to be"), 1e2)

func TestHash(t *testing.T) {
	h := New()
	h.Write([]byte("to be or not to be"))
	t.Log(h.Sum64(), 1988685042348123509)
	h = New()
	h.Write2([]byte("to be or not to be"))
	t.Log(h.Sum64(), 1988685042348123509)
	h = New()
	h.Write3([]byte("to be or not to be"))
	t.Log(h.Sum64(), 1988685042348123509)
}

var sink uint64

func BenchmarkXX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = xxhash.Checksum64(testBuff)
	}
}

func BenchmarkSH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := New()
		h.Write(testBuff)
		sink = h.Sum64()
	}
}

func BenchmarkSH2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := New()
		h.Write2(testBuff)
		sink = h.Sum64()
	}
}

func BenchmarkSH3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := New()
		h.Write3(testBuff)
		sink = h.Sum64()
	}
}

func TestDiffuse(t *testing.T) {
	for _, tst := range [...][2]uint64{{94203824938, 17289265692384716055}, {0xDEADBEEF, 12110756357096144265}, {0, 0},
		{1, 15197155197312260123}, {2, 1571904453004118546}, {3, 16467633989910088880}} {
		if v := diffuse(tst[0]); v != tst[1] {
			diffuseTest(t, tst[0], tst[1])
		}
	}
}

func diffuseTest(t *testing.T, x, y uint64) {
	/*
	   assert_eq!(diffuse(x), y);
	         assert_eq!(x, undiffuse(y));
	         assert_eq!(undiffuse(diffuse(x)), x);
	*/
	if v := diffuse(x); v != y {
		t.Fatalf("diffuse(0x%x) != 0x%x (got: 0x%x)", x, y, v)
	}

	if v := undiffuse(y); v != x {
		t.Fatalf("0x%x != undiffuse(0x%x) (got: 0x%x)", x, y, v)
	}
	if v := undiffuse(diffuse(x)); v != x {
		t.Fatalf("undiffuse(diffuse(0x%x)) != 0x%x (got: 0x%x)", x, y, v)
	}
}
