package sotautils

import (
	"testing"
)

func TestGetBandFromFrequency(t *testing.T) {
	tests := []struct {
		name     string
		freqKHz  float64
		expected string
	}{
		{"80m band", 3600, "80m"},
		{"40m band", 7100, "40m"},
		{"20m band", 14200, "20m"},
		{"2m band", 145000, "2m"},
		{"70cm band", 435000, "70cm"},
		{"Out of band", 5000, ""},
		{"160m band", 1900, "160m"},
		{"10m band", 28500, "10m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBandFromFrequency(tt.freqKHz)
			if result != tt.expected {
				t.Errorf("GetBandFromFrequency(%v) = %v, want %v", tt.freqKHz, result, tt.expected)
			}
		})
	}
}

func TestGetBandFromFrequencyMHz(t *testing.T) {
	tests := []struct {
		name     string
		freqMHz  float64
		expected string
	}{
		{"14 MHz in 20m", 14.200, "20m"},
		{"144 MHz in 2m", 145.500, "2m"},
		{"432 MHz in 70cm", 433.500, "70cm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBandFromFrequencyMHz(tt.freqMHz)
			if result != tt.expected {
				t.Errorf("GetBandFromFrequencyMHz(%v) = %v, want %v", tt.freqMHz, result, tt.expected)
			}
		})
	}
}

func TestIsFrequencyInBand(t *testing.T) {
	tests := []struct {
		name     string
		freqKHz  float64
		bandName string
		expected bool
	}{
		{"7100 kHz in 40m", 7100, "40m", true},
		{"14200 kHz in 20m", 14200, "20m", true},
		{"7100 kHz not in 20m", 7100, "20m", false},
		{"Invalid band", 7100, "invalid", false},
		{"145000 kHz in 2m", 145000, "2m", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFrequencyInBand(tt.freqKHz, tt.bandName)
			if result != tt.expected {
				t.Errorf("IsFrequencyInBand(%v, %v) = %v, want %v", tt.freqKHz, tt.bandName, result, tt.expected)
			}
		})
	}
}

func TestGetAllBandNames(t *testing.T) {
	names := GetAllBandNames()
	if len(names) == 0 {
		t.Error("GetAllBandNames() returned empty slice")
	}

	// Check that some expected bands are present
	expectedBands := []string{"80m", "40m", "20m", "2m", "70cm"}
	for _, expected := range expectedBands {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected band %s not found in band names", expected)
		}
	}
}

func TestRSGBBandPlanCompleteness(t *testing.T) {
	// Verify all expected bands are present
	expectedBands := []string{
		"2190m", "630m", "160m",
		"80m", "60m", "40m", "30m", "20m", "17m", "15m", "12m", "10m",
		"6m", "4m", "2m", "70cm",
		"23cm", "13cm", "9cm", "6cm", "3cm",
	}

	for _, band := range expectedBands {
		if _, exists := RSGBBandPlan[band]; !exists {
			t.Errorf("Expected band %s not found in RSGBBandPlan", band)
		}
	}
}

func TestBandRangeValidity(t *testing.T) {
	// Verify that all bands have valid frequency ranges (lower < upper)
	for name, band := range RSGBBandPlan {
		if band.LowerFreq >= band.UpperFreq {
			t.Errorf("Band %s has invalid range: %v - %v kHz", name, band.LowerFreq, band.UpperFreq)
		}
		if band.Name != name {
			t.Errorf("Band %s has mismatched name field: %s", name, band.Name)
		}
	}
}
