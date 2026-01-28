package math

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func TestAbsolute(t *testing.T) {
	// We only test the positive case
	result := Absolute(5)
	expected := 5

	if result != expected {
		t.Errorf("Absolute(5) = %d; want %d", result, expected)
	}

	// MISSING: We forgot to add a test case for Absolute(-5)
}
