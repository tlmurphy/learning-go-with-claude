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

func TestSafeIntToInt8(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int8
		expectOK bool
	}{
		{"zero", 0, 0, true},
		{"positive in range", 42, 42, true},
		{"negative in range", -100, -100, true},
		{"max int8", 127, 127, true},
		{"min int8", -128, -128, true},
		{"overflow positive", 128, 0, false},
		{"overflow negative", -129, 0, false},
		{"large positive", 1000, 0, false},
		{"large negative", -1000, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := SafeIntToInt8(tt.input)
			if ok != tt.expectOK {
				t.Errorf("SafeIntToInt8(%d): expected ok=%t, got ok=%t.\n"+
					"  Hint: Check if the value is within int8 range (-128 to 127) before converting.",
					tt.input, tt.expectOK, ok)
			}
			if got != tt.expected {
				t.Errorf("SafeIntToInt8(%d): expected value=%d, got value=%d.\n"+
					"  Hint: Return int8(n) if in range, or 0 if out of range.",
					tt.input, tt.expected, got)
			}
		})
	}
}

func TestStringByteRuneAnalysis(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBytes int
		wantRunes int
		wantFirst rune
		wantLast  rune
	}{
		{
			name:      "ASCII only",
			input:     "hello",
			wantBytes: 5, wantRunes: 5,
			wantFirst: 'h', wantLast: 'o',
		},
		{
			name:      "with Chinese characters",
			input:     "Go世界",
			wantBytes: 8, wantRunes: 4,
			wantFirst: 'G', wantLast: '界',
		},
		{
			name:      "emoji",
			input:     "🎉🌍",
			wantBytes: 8, wantRunes: 2,
			wantFirst: '🎉', wantLast: '🌍',
		},
		{
			name:      "single character",
			input:     "X",
			wantBytes: 1, wantRunes: 1,
			wantFirst: 'X', wantLast: 'X',
		},
		{
			name:      "empty string",
			input:     "",
			wantBytes: 0, wantRunes: 0,
			wantFirst: 0, wantLast: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, gotRunes, gotFirst, gotLast := StringByteRuneAnalysis(tt.input)

			if gotBytes != tt.wantBytes {
				t.Errorf("StringByteRuneAnalysis(%q) byteCount = %d, want %d.\n"+
					"  Hint: Use len(s) for byte count.",
					tt.input, gotBytes, tt.wantBytes)
			}
			if gotRunes != tt.wantRunes {
				t.Errorf("StringByteRuneAnalysis(%q) runeCount = %d, want %d.\n"+
					"  Hint: Convert to []rune first, then use len().",
					tt.input, gotRunes, tt.wantRunes)
			}
			if gotFirst != tt.wantFirst {
				t.Errorf("StringByteRuneAnalysis(%q) firstRune = %q, want %q.\n"+
					"  Hint: After converting to []rune, the first rune is at index 0.",
					tt.input, gotFirst, tt.wantFirst)
			}
			if gotLast != tt.wantLast {
				t.Errorf("StringByteRuneAnalysis(%q) lastRune = %q, want %q.\n"+
					"  Hint: After converting to []rune, the last rune is at index len(runes)-1.",
					tt.input, gotLast, tt.wantLast)
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

func TestEvaluateScore(t *testing.T) {
	tests := []struct {
		name           string
		score          UserScore
		wantPercentage float64
		wantRating     string
		wantPassing    bool
	}{
		{"perfect score", UserScore(200), 100.0, "excellent", true},
		{"excellent", UserScore(190), 95.0, "excellent", true},
		{"good boundary", UserScore(140), 70.0, "good", true},
		{"good", UserScore(150), 75.0, "good", true},
		{"fair boundary", UserScore(100), 50.0, "fair", true},
		{"fair", UserScore(110), 55.0, "fair", true},
		{"needs improvement", UserScore(80), 40.0, "needs improvement", false},
		{"zero score", UserScore(0), 0.0, "needs improvement", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPct, gotRating, gotPassing := EvaluateScore(tt.score)

			if diff := gotPct - tt.wantPercentage; diff > 0.01 || diff < -0.01 {
				t.Errorf("EvaluateScore(%d) percentage = %.2f, want %.2f.\n"+
					"  Hint: percentage = float64(score) / float64(MaxScore) * 100",
					tt.score, gotPct, tt.wantPercentage)
			}
			if gotRating != tt.wantRating {
				t.Errorf("EvaluateScore(%d) rating = %q, want %q.\n"+
					"  Hint: Check thresholds: >=90 excellent, >=70 good, >=50 fair, else needs improvement.",
					tt.score, gotRating, tt.wantRating)
			}
			if gotPassing != tt.wantPassing {
				t.Errorf("EvaluateScore(%d) passing = %t, want %t.\n"+
					"  Hint: A score is passing if the percentage >= 50.",
					tt.score, gotPassing, tt.wantPassing)
			}
		})
	}
}
