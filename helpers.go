package seahash

const (
	pcgRNG    = 0x6eed0e9da4d94a4f
	pcgRNGrev = 0x2f72b4215a3d8caf
)

func diffuse(x uint64) uint64 {
	x *= pcgRNG
	a, b := x>>32, x>>60
	x ^= a >> b
	x *= pcgRNG
	return x
}

func undiffuse(x uint64) uint64 {
	x *= pcgRNGrev
	a, b := x>>32, x>>60
	x ^= a >> b
	x *= pcgRNGrev
	return x
}
