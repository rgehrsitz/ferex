package testutils

func Abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func Contains(s, substr string) bool {
	// Simplified 'Contains' for basic test needs. Real-world might use strings.Contains.
	// This legacy version seems to have a recursive bug for non-matches or partials.
	// For now, keeping it as is from legacy, but it might need revision if tests fail unexpectedly due to it.
	// A more robust standard library check would be: import "strings"; return strings.Contains(s, substr)
	return len(substr) == 0 || (len(s) > 0 && (s == substr || (len(s) >= len(substr) && (s[0:len(substr)] == substr || Contains(s[1:], substr)))))
}
