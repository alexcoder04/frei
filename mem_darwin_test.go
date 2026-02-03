//go:build darwin

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

	// MemBuffers is always 0 on macOS, but should be non-negative
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

func TestGetVMStats(t *testing.T) {
	stats, err := getVMStats()
	if err != nil {
		t.Fatalf("getVMStats() returned error: %v", err)
	}

	// Should have some basic stats
	requiredKeys := []string{
		"Pages free",
		"Pages active",
		"Pages inactive",
		"Pages wired down",
	}

	for _, key := range requiredKeys {
		if _, ok := stats[key]; !ok {
			t.Errorf("getVMStats() missing expected key: %s", key)
		}
	}

	// All values should be non-negative
	for key, value := range stats {
		if value < 0 {
			t.Errorf("getVMStats()[%s] should be non-negative, got %f", key, value)
		}
	}
}

func TestMemoryConsistency(t *testing.T) {
	data, err := GetMemInfo()
	if err != nil {
		t.Fatalf("GetMemInfo() returned error: %v", err)
	}

	// Basic sanity check: used + free + cached should be close to total
	// (allowing for some variance due to timing and measurement methods)
	accounted := data.MemUsed + data.MemFree + data.MemCached

	// Allow 20% variance due to different measurement methods
	variance := data.MemTotal * 0.2
	if accounted < data.MemTotal-variance || accounted > data.MemTotal+variance {
		t.Logf("Memory accounting may be off: Used(%f) + Free(%f) + Cached(%f) = %f, Total = %f",
			data.MemUsed, data.MemFree, data.MemCached, accounted, data.MemTotal)
	}
}
