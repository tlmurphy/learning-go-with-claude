package variables

import (
	"testing"
)

func TestDeclareVariables(t *testing.T) {
	i, s, b := DeclareVariables()

	if i != 42 {
		t.Errorf("Expected int value 42, got %d. Declare a variable with value 42 and return it.", i)
	}
	if s != "hello" {
		t.Errorf("Expected string value %q, got %q. Declare a string variable with value \"hello\".", "hello", s)
	}
	if b != true {
		t.Errorf("Expected bool value true, got %t. Declare a bool variable with value true.", b)
	}
}

func TestZeroValues(t *testing.T) {
	i, f, s, b := ZeroValues()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
		hint     string
	}{
		{"int", i, 0, "The zero value for int is 0. Use `var x int` without assignment."},
		{"float64", f, 0.0, "The zero value for float64 is 0.0. Use `var x float64` without assignment."},
		{"string", s, "", "The zero value for string is \"\" (empty string). Use `var x string` without assignment."},
		{"bool", b, false, "The zero value for bool is false. Use `var x bool` without assignment."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Zero value for %s: expected %v, got %v.\n  Hint: %s",
					tt.name, tt.expected, tt.got, tt.hint)
			}
		})
	}
}

func TestGetStatusCategories(t *testing.T) {
	info, success, redirect, clientErr, serverErr := GetStatusCategories()

	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"StatusInfo", info, 1},
		{"StatusSuccess", success, 2},
		{"StatusRedirect", redirect, 3},
		{"StatusClientError", clientErr, 4},
		{"StatusServerError", serverErr, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s: expected %d, got %d.\n"+
					"  Hint: Use iota in a const block starting at 1. Remember that iota\n"+
					"  starts at 0, so you need `iota + 1` or a blank identifier `_` to skip 0.",
					tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		input    Celsius
		expected Kelvin
	}{
		{"boiling water", Celsius(100), Kelvin(373.15)},
		{"freezing water", Celsius(0), Kelvin(273.15)},
		{"absolute zero", Celsius(-273.15), Kelvin(0)},
		{"body temperature", Celsius(37), Kelvin(310.15)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CelsiusToKelvin(tt.input)
			if diff := float64(got) - float64(tt.expected); diff > 0.001 || diff < -0.001 {
				t.Errorf("CelsiusToKelvin(%v) = %v, want %v.\n"+
					"  Hint: Kelvin = Celsius + 273.15. Remember to use type conversions!",
					tt.input, got, tt.expected)
			}
		})
	}
}

func TestKelvinToCelsius(t *testing.T) {
	tests := []struct {
		name     string
		input    Kelvin
		expected Celsius
	}{
		{"boiling water", Kelvin(373.15), Celsius(100)},
		{"freezing water", Kelvin(273.15), Celsius(0)},
		{"absolute zero", Kelvin(0), Celsius(-273.15)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KelvinToCelsius(tt.input)
			if diff := float64(got) - float64(tt.expected); diff > 0.001 || diff < -0.001 {
				t.Errorf("KelvinToCelsius(%v) = %v, want %v.\n"+
					"  Hint: Celsius = Kelvin - 273.15. Remember to use type conversions!",
					tt.input, got, tt.expected)
			}
		})
	}
}

func TestAbsoluteZeroCelsius(t *testing.T) {
	got := AbsoluteZeroCelsius()
	expected := Celsius(-273.15)
	if diff := float64(got) - float64(expected); diff > 0.001 || diff < -0.001 {
		t.Errorf("AbsoluteZeroCelsius() = %v, want %v.\n"+
			"  Hint: Convert Kelvin(0) to Celsius using your KelvinToCelsius function.",
			got, expected)
	}
}

func TestSwapInts(t *testing.T) {
	tests := []struct {
		name  string
		a, b  int
		wantA int
		wantB int
	}{
		{"positive numbers", 1, 2, 2, 1},
		{"with zero", 0, 42, 42, 0},
		{"same value", 5, 5, 5, 5},
		{"negative numbers", -3, -7, -7, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := SwapInts(tt.a, tt.b)
			if gotA != tt.wantA || gotB != tt.wantB {
				t.Errorf("SwapInts(%d, %d) = (%d, %d), want (%d, %d).\n"+
					"  Hint: Go supports multiple assignment: a, b = b, a",
					tt.a, tt.b, gotA, gotB, tt.wantA, tt.wantB)
			}
		})
	}
}

func TestSwapStrings(t *testing.T) {
	tests := []struct {
		name  string
		a, b  string
		wantA string
		wantB string
	}{
		{"simple", "hello", "world", "world", "hello"},
		{"with empty", "", "test", "test", ""},
		{"same value", "go", "go", "go", "go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := SwapStrings(tt.a, tt.b)
			if gotA != tt.wantA || gotB != tt.wantB {
				t.Errorf("SwapStrings(%q, %q) = (%q, %q), want (%q, %q).\n"+
					"  Hint: Same pattern as SwapInts — Go's multiple return works with any type.",
					tt.a, tt.b, gotA, gotB, tt.wantA, tt.wantB)
			}
		})
	}
}
