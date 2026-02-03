//go:build linux

package main

import (
	"testing"
)

func TestGetMemInfo(t *testing.T) {
	data, err := GetMemInfo()
	if err != nil {
		t.Fatalf("GetMemInfo() returned error: %v", err)
	}

	// MemTotal should be positive (system has memory)
	if data.MemTotal <= 0 {
		t.Errorf("MemTotal should be positive, got %f", data.MemTotal)
	}

	// MemUsed should be positive (some memory is always in use)
	if data.MemUsed <= 0 {
		t.Errorf("MemUsed should be positive, got %f", data.MemUsed)
	}

	// MemUsed should be less than MemTotal
	if data.MemUsed >= data.MemTotal {
		t.Errorf("MemUsed (%f) should be less than MemTotal (%f)", data.MemUsed, data.MemTotal)
	}

	// MemFree should be non-negative
	if data.MemFree < 0 {
		t.Errorf("MemFree should be non-negative, got %f", data.MemFree)
	}

	// MemAvailable should be positive and less than or equal to MemTotal
	if data.MemAvailable <= 0 {
		t.Errorf("MemAvailable should be positive, got %f", data.MemAvailable)
	}
	if data.MemAvailable > data.MemTotal {
		t.Errorf("MemAvailable (%f) should not exceed MemTotal (%f)", data.MemAvailable, data.MemTotal)
	}

	// MemCached should be non-negative
	if data.MemCached < 0 {
		t.Errorf("MemCached should be non-negative, got %f", data.MemCached)
	}

	// MemBuffers should be non-negative
	if data.MemBuffers < 0 {
		t.Errorf("MemBuffers should be non-negative, got %f", data.MemBuffers)
	}

	// SwapTotal should be non-negative (swap may be disabled)
	if data.SwapTotal < 0 {
		t.Errorf("SwapTotal should be non-negative, got %f", data.SwapTotal)
	}

	// If there's swap, check swap values
	if data.SwapTotal > 0 {
		if data.SwapFree < 0 {
			t.Errorf("SwapFree should be non-negative, got %f", data.SwapFree)
		}
		if data.SwapUsed < 0 {
			t.Errorf("SwapUsed should be non-negative, got %f", data.SwapUsed)
		}
		if data.SwapFree > data.SwapTotal {
			t.Errorf("SwapFree (%f) should not exceed SwapTotal (%f)", data.SwapFree, data.SwapTotal)
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		input    string
		wantKey  string
		wantVal  float64
	}{
		{"MemTotal:       16384000 kB", "MemTotal", 16384000},
		{"MemFree:         1234567 kB", "MemFree", 1234567},
		{"Buffers:              0 kB", "Buffers", 0},
		{"SwapTotal:      2097152 kB", "SwapTotal", 2097152},
	}

	for _, tt := range tests {
		key, val := parseLine(tt.input)
		if key != tt.wantKey {
			t.Errorf("parseLine(%q) key = %q, want %q", tt.input, key, tt.wantKey)
		}
		if val != tt.wantVal {
			t.Errorf("parseLine(%q) value = %f, want %f", tt.input, val, tt.wantVal)
		}
	}
}

func TestMemoryConsistency(t *testing.T) {
	data, err := GetMemInfo()
	if err != nil {
		t.Fatalf("GetMemInfo() returned error: %v", err)
	}

	// Basic sanity check: used + free + cached + buffers should be close to total
	accounted := data.MemUsed + data.MemFree + data.MemCached + data.MemBuffers

	// Allow 20% variance
	variance := data.MemTotal * 0.2
	if accounted < data.MemTotal-variance || accounted > data.MemTotal+variance {
		t.Logf("Memory accounting may be off: Used(%f) + Free(%f) + Cached(%f) + Buffers(%f) = %f, Total = %f",
			data.MemUsed, data.MemFree, data.MemCached, data.MemBuffers, accounted, data.MemTotal)
	}
}
