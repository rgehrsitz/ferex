package testutils

func Abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func Contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) > 0 && (s == substr || (len(s) >= len(substr) && (s[0:len(substr)] == substr || Contains(s[1:], substr)))))
}
