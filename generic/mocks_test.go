package generic

var (
	nilInt   *int64
	nilUint  *uint64
	nilFloat *float64
)

func uintPointer(u uint) *uint {
	return &u
}

func intPointer(i int) *int {
	return &i
}

func floatPointer(f float64) *float64 {
	return &f
}

func stringPointer(s string) *string {
	return &s
}

type mapKey struct {
	key string
}
