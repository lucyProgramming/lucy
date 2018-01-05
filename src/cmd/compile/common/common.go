package common

var (
	E = 0.0001
)

func Float32Equal(f float32, f2 float32) bool {
	return Float64Equal(float64(f), float64(f2))
}
func Float64Equal(f float64, f2 float64) bool {
	return f >= f2-E && f <= f2+E
}
