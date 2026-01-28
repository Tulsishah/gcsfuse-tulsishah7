package math

// Add takes two integers and returns their sum
func Add(a, b int) int {
	return a + b
}

func Absolute(n int) int {
	if n < 0 {
		// This is the "negative path"
		return -n
	}
	// This is the "positive path"
	return n
}
