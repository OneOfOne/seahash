package seahash

import "encoding/binary"
import "unsafe"

// based on https://github.com/ticki/tfs/blob/master/seahash/src/reference.rs

// https://github.com/ticki/tfs/blob/master/seahash/src/buffer.rs
var encLE = binary.LittleEndian

func New() *seahash {
	return NewWithSeeds(0x16f11fe89b0d677c, 0xb480a793d8e6c86c, 0x6fe2e5aaf078ebc9, 0x14f994a4c5259381)
}

func NewWithSeeds(k1, k2, k3, k4 uint64) *seahash {
	return &seahash{a: k1, b: k2, c: k3, d: k4}
}

type seahash struct {
	a, b, c, d uint64
	len        int
}

func (s *seahash) writeUint64(x uint64) {
	a := diffuse(s.a ^ x)
	s.a, s.b, s.c, s.d = s.b, s.c, s.d, a
}

func (s *seahash) Write(p []byte) (int64, error) {
	var i int
	for ; i < len(p)-7; i += 8 {
		s.writeUint64(encLE.Uint64(p[i:]))
	}
	for ; i < len(p)-3; i += 4 {
		s.writeUint64(uint64(encLE.Uint32(p[i:])))
	}
	for ; i < len(p)-1; i += 2 {
		s.writeUint64(uint64(encLE.Uint16(p[i:])))
	}
	for ; i < len(p); i++ {
		s.writeUint64(uint64(p[i]))
	}
	s.len += len(p)
	return int64(len(p)), nil
}

func (s *seahash) Write2(p []byte) (int64, error) {
	var i int
	for ; i < len(p)-32; i += 32 {
		s.a = diffuse(s.a ^ encLE.Uint64(p[i:]))
		s.b = diffuse(s.b ^ encLE.Uint64(p[i+8:]))
		s.c = diffuse(s.c ^ encLE.Uint64(p[i+16:]))
		s.d = diffuse(s.d ^ encLE.Uint64(p[i+24:]))
	}
	for ; i < len(p)-7; i += 8 {
		s.writeUint64(encLE.Uint64(p[i:]))
	}
	for ; i < len(p)-3; i += 4 {
		s.writeUint64(uint64(encLE.Uint32(p[i:])))
	}
	for ; i < len(p)-1; i += 2 {
		s.writeUint64(uint64(encLE.Uint16(p[i:])))
	}
	for ; i < len(p); i++ {
		s.writeUint64(uint64(p[i]))
	}
	s.len += len(p)
	return int64(len(p)), nil
}

func (s *seahash) Write3(p []byte) (int64, error) {
	var i int
	for ; i < len(p)-32; i += 32 {
		pp := unsafe.Pointer(&p[i])
		s.a = diffuse(s.a ^ unsafeUint64(pp, 0))
		s.b = diffuse(s.b ^ unsafeUint64(pp, 8))
		s.c = diffuse(s.c ^ unsafeUint64(pp, 16))
		s.d = diffuse(s.d ^ unsafeUint64(pp, 24))
	}
	for ; i < len(p)-7; i += 8 {
		pp := unsafe.Pointer(&p[i])
		s.writeUint64(unsafeUint64(pp, 0))
	}
	for ; i < len(p)-3; i += 4 {
		s.writeUint64(uint64(encLE.Uint32(p[i:])))
	}
	for ; i < len(p)-1; i += 2 {
		s.writeUint64(uint64(encLE.Uint16(p[i:])))
	}
	for ; i < len(p); i++ {
		s.writeUint64(uint64(p[i]))
	}
	s.len += len(p)
	return int64(len(p)), nil
}

func unsafeUint64(p unsafe.Pointer, idx uintptr) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(p) + idx))
}

func (s *seahash) Sum64() uint64 {
	return diffuse(s.a ^ s.b ^ s.c ^ s.d ^ uint64(s.len))
}
