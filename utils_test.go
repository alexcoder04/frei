package main

import (
	"testing"
)

func TestToFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"123", 123},
		{"123.45", 123.45},
		{"0", 0},
		{"", 0},
		{"invalid", 0},
		{"16384000", 16384000},
	}

	for _, tt := range tests {
		got := toFloat(tt.input)
		if got != tt.want {
			t.Errorf("toFloat(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}

func TestToHumanStr(t *testing.T) {
	tests := []struct {
		value float64
		human bool
		want  string
	}{
		{1024, false, "1 M"},       // 1024 KB = 1 MB when not human
		{2048, false, "2 M"},       // 2048 KB = 2 MB when not human
		{1024, true, "1.00 M"},     // 1024 KB = 1 MB in human readable
		{512, true, "512.00 K"},    // 512 KB stays in KB
		{1048576, true, "1.00 G"},  // 1 GB in human readable
	}

	for _, tt := range tests {
		got := toHumanStr(tt.value, tt.human)
		if got != tt.want {
			t.Errorf("toHumanStr(%f, %v) = %q, want %q", tt.value, tt.human, got, tt.want)
		}
	}
}

func TestGetTerminalWidth(t *testing.T) {
	width := getTerminalWidth()
	// Should return a reasonable value (fallback is 80)
	if width <= 0 {
		t.Errorf("getTerminalWidth() = %d, want positive value", width)
	}
	// Should be at least the fallback value or more if terminal is wider
	// In a test environment, we might get 80 (fallback)
	if width < 10 || width > 1000 {
		t.Errorf("getTerminalWidth() = %d, seems unreasonable", width)
	}
}
