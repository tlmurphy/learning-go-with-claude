package controlflow

import (
	"reflect"
	"testing"
)

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

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected []string
	}{
		{
			name:     "up to 1",
			n:        1,
			expected: []string{"1"},
		},
		{
			name:     "up to 5",
			n:        5,
			expected: []string{"1", "2", "Fizz", "4", "Buzz"},
		},
		{
			name: "up to 15",
			n:    15,
			expected: []string{
				"1", "2", "Fizz", "4", "Buzz",
				"Fizz", "7", "8", "Fizz", "Buzz",
				"11", "Fizz", "13", "14", "FizzBuzz",
			},
		},
		{
			name:     "zero",
			n:        0,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FizzBuzz(tt.n)
			if got == nil && tt.n >= 0 {
				t.Fatalf("FizzBuzz(%d) returned nil. Return an empty slice for n=0, not nil.", tt.n)
			}
			if len(got) != len(tt.expected) {
				t.Fatalf("FizzBuzz(%d) returned %d elements, want %d.\n  got:  %v\n  want: %v",
					tt.n, len(got), len(tt.expected), got, tt.expected)
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("FizzBuzz(%d)[%d] = %q, want %q.\n"+
						"  Hint: Check your divisibility conditions. Test divisible-by-15 BEFORE 3 and 5.",
						tt.n, i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestClassifyAge(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected string
	}{
		{"negative", -1, "invalid"},
		{"baby", 0, "child"},
		{"child", 8, "child"},
		{"child boundary", 12, "child"},
		{"teenager start", 13, "teenager"},
		{"teenager", 15, "teenager"},
		{"teenager boundary", 17, "teenager"},
		{"adult start", 18, "adult"},
		{"adult", 30, "adult"},
		{"adult boundary", 64, "adult"},
		{"senior start", 65, "senior"},
		{"senior", 80, "senior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyAge(tt.age)
			if got != tt.expected {
				t.Errorf("ClassifyAge(%d) = %q, want %q.\n"+
					"  Hint: Use if/else if chain with the boundaries: <0, <13, <18, <65, >=65.",
					tt.age, got, tt.expected)
			}
		})
	}
}

func TestDayType(t *testing.T) {
	tests := []struct {
		name     string
		day      string
		expected string
	}{
		{"monday lowercase", "monday", "weekday"},
		{"Monday mixed case", "Monday", "weekday"},
		{"FRIDAY uppercase", "FRIDAY", "weekday"},
		{"saturday", "saturday", "weekend"},
		{"SUNDAY", "SUNDAY", "weekend"},
		{"Wednesday", "Wednesday", "weekday"},
		{"invalid", "Funday", "invalid"},
		{"empty", "", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DayType(tt.day)
			if got != tt.expected {
				t.Errorf("DayType(%q) = %q, want %q.\n"+
					"  Hint: Normalize the input to lowercase before switching. Use strings.ToLower().",
					tt.day, got, tt.expected)
			}
		})
	}
}

func TestCountUnicodeCategories(t *testing.T) {
	tests := []struct {
		name                                                    string
		input                                                   string
		wantUpper, wantLower, wantDigits, wantSpaces, wantOther int
	}{
		{
			name:      "mixed ASCII",
			input:     "Hello World 123!",
			wantUpper: 2, wantLower: 8, wantDigits: 3, wantSpaces: 2, wantOther: 1,
		},
		{
			name:      "all lowercase",
			input:     "abc",
			wantUpper: 0, wantLower: 3, wantDigits: 0, wantSpaces: 0, wantOther: 0,
		},
		{
			name:      "password check",
			input:     "P@ssw0rd!",
			wantUpper: 1, wantLower: 4, wantDigits: 1, wantSpaces: 0, wantOther: 3,
		},
		{
			name:      "empty string",
			input:     "",
			wantUpper: 0, wantLower: 0, wantDigits: 0, wantSpaces: 0, wantOther: 0,
		},
		{
			name:      "whitespace types",
			input:     "a\tb\nc",
			wantUpper: 0, wantLower: 3, wantDigits: 0, wantSpaces: 2, wantOther: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, l, d, s, o := CountUnicodeCategories(tt.input)
			if u != tt.wantUpper {
				t.Errorf("upper: got %d, want %d.\n  Hint: Use unicode.IsUpper(r) in a tagless switch.", u, tt.wantUpper)
			}
			if l != tt.wantLower {
				t.Errorf("lower: got %d, want %d.\n  Hint: Use unicode.IsLower(r).", l, tt.wantLower)
			}
			if d != tt.wantDigits {
				t.Errorf("digits: got %d, want %d.\n  Hint: Use unicode.IsDigit(r).", d, tt.wantDigits)
			}
			if s != tt.wantSpaces {
				t.Errorf("spaces: got %d, want %d.\n  Hint: Use unicode.IsSpace(r). This includes tabs and newlines.", s, tt.wantSpaces)
			}
			if o != tt.wantOther {
				t.Errorf("other: got %d, want %d.\n  Hint: The default case in your switch catches everything else.", o, tt.wantOther)
			}
		})
	}
}

func TestFindInMatrix(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	tests := []struct {
		name      string
		matrix    [][]int
		target    int
		wantRow   int
		wantCol   int
		wantFound bool
	}{
		{"find 1 (top-left)", matrix, 1, 0, 0, true},
		{"find 5 (center)", matrix, 5, 1, 1, true},
		{"find 9 (bottom-right)", matrix, 9, 2, 2, true},
		{"find 6", matrix, 6, 1, 2, true},
		{"not found", matrix, 42, -1, -1, false},
		{"empty matrix", [][]int{}, 1, -1, -1, false},
		{"single element found", [][]int{{7}}, 7, 0, 0, true},
		{"single element not found", [][]int{{7}}, 3, -1, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, found := FindInMatrix(tt.matrix, tt.target)
			if found != tt.wantFound || row != tt.wantRow || col != tt.wantCol {
				t.Errorf("FindInMatrix(matrix, %d) = (%d, %d, %t), want (%d, %d, %t).\n"+
					"  Hint: Use a labeled break to exit both loops when you find the target.",
					tt.target, row, col, found, tt.wantRow, tt.wantCol, tt.wantFound)
			}
		})
	}
}

func TestStateMachine(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		expected string
	}{
		{"no commands", []string{}, "idle"},
		{"start", []string{"start"}, "running"},
		{"start and stop", []string{"start", "stop"}, "stopped"},
		{"start, pause, start", []string{"start", "pause", "start"}, "running"},
		{"full cycle", []string{"start", "stop", "reset"}, "idle"},
		{"invalid from idle", []string{"stop"}, "idle"},
		{"invalid from running", []string{"start", "reset"}, "running"},
		{"complex sequence", []string{"start", "pause", "start", "stop", "reset", "start"}, "running"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StateMachine(tt.commands)
			if got != tt.expected {
				t.Errorf("StateMachine(%v) = %q, want %q.\n"+
					"  Hint: Use a switch on the current state, then check the command.\n"+
					"  Valid transitions: idle+start->running, running+stop->stopped,\n"+
					"  running+pause->idle, stopped+reset->idle.",
					tt.commands, got, tt.expected)
			}
		})
	}
}

func TestCollatzSteps(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"n=1", 1, 0},
		{"n=2", 2, 1},
		{"n=6", 6, 8},
		{"n=11", 11, 14},
		{"n=27", 27, 111},
		{"n=0 (invalid)", 0, -1},
		{"n=-5 (invalid)", -5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CollatzSteps(tt.n)
			if got != tt.expected {
				t.Errorf("CollatzSteps(%d) = %d, want %d.\n"+
					"  Hint: Loop while n != 1. If even, n = n/2. If odd, n = 3*n + 1. Count each step.",
					tt.n, got, tt.expected)
			}
		})
	}
}

func TestProcessRecords(t *testing.T) {
	tests := []struct {
		name        string
		records     []string
		wantResults map[string]string
		wantErrors  []string
	}{
		{
			name:        "basic key-value pairs",
			records:     []string{"name:Alice", "age:30"},
			wantResults: map[string]string{"name": "Alice", "age": "30"},
			wantErrors:  []string{},
		},
		{
			name:        "skip empty strings",
			records:     []string{"", "name:Bob", ""},
			wantResults: map[string]string{"name": "Bob"},
			wantErrors:  []string{},
		},
		{
			name:        "missing colon goes to errors",
			records:     []string{"name:Alice", "bad-record", "age:30"},
			wantResults: map[string]string{"name": "Alice", "age": "30"},
			wantErrors:  []string{"bad-record"},
		},
		{
			name:        "STOP halts processing",
			records:     []string{"name:Alice", "STOP:now", "age:30"},
			wantResults: map[string]string{"name": "Alice"},
			wantErrors:  []string{},
		},
		{
			name:        "empty key goes to errors",
			records:     []string{":value"},
			wantResults: map[string]string{},
			wantErrors:  []string{":value"},
		},
		{
			name:        "empty value gets default",
			records:     []string{"theme:"},
			wantResults: map[string]string{"theme": "default"},
			wantErrors:  []string{},
		},
		{
			name:        "combined scenario",
			records:     []string{"", "host:localhost", "bad", ":oops", "port:", "STOP:", "never:seen"},
			wantResults: map[string]string{"host": "localhost", "port": "default"},
			wantErrors:  []string{"bad", ":oops"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, gotErrors := ProcessRecords(tt.records)

			if gotResults == nil {
				t.Fatal("ProcessRecords returned nil results map. Initialize your map with make().")
			}
			if gotErrors == nil {
				t.Fatal("ProcessRecords returned nil errors slice. Return an empty slice, not nil.")
			}

			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("results = %v, want %v.\n"+
					"  Hint: Use strings.Index or strings.SplitN to split on the first ':'.",
					gotResults, tt.wantResults)
			}
			if !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("errors = %v, want %v.\n"+
					"  Hint: Records without ':' and records with empty keys go to errors.",
					gotErrors, tt.wantErrors)
			}
		})
	}
}
