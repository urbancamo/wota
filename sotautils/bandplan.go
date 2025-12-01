package sotautils

// BandRange represents the frequency range for an amateur radio band
type BandRange struct {
	Name      string  // Common name (e.g., "80m", "2m")
	LowerFreq float64 // Lower frequency in kHz
	UpperFreq float64 // Upper frequency in kHz
}

// RSGBBandPlan contains UK amateur radio band allocations based on RSGB 2024/2025 band plans
// Sources:
// - https://rsgb.org/main/operating/band-plans/
// - https://rsgb.org/main/operating/band-plans/hf/
// - https://rsgb.org/main/operating/band-plans/lf-mf-bands/
// - https://rsgb.org/main/operating/band-plans/vhf-uhf/
var RSGBBandPlan = map[string]BandRange{
	// LF Bands
	"2190m": {Name: "2190m", LowerFreq: 135.7, UpperFreq: 137.8},
	"630m":  {Name: "630m", LowerFreq: 472, UpperFreq: 479},

	// MF Band
	"160m": {Name: "160m", LowerFreq: 1810, UpperFreq: 2000},

	// HF Bands
	"80m": {Name: "80m", LowerFreq: 3500, UpperFreq: 3800},
	"60m": {Name: "60m", LowerFreq: 5258.5, UpperFreq: 5406.5}, // Multiple channels, approximate range
	"40m": {Name: "40m", LowerFreq: 7000, UpperFreq: 7200},
	"30m": {Name: "30m", LowerFreq: 10100, UpperFreq: 10150},
	"20m": {Name: "20m", LowerFreq: 14000, UpperFreq: 14350},
	"17m": {Name: "17m", LowerFreq: 18068, UpperFreq: 18168},
	"15m": {Name: "15m", LowerFreq: 21000, UpperFreq: 21450},
	"12m": {Name: "12m", LowerFreq: 24890, UpperFreq: 24990},
	"10m": {Name: "10m", LowerFreq: 28000, UpperFreq: 29700},

	// VHF Bands
	"6m": {Name: "6m", LowerFreq: 50000, UpperFreq: 52000},
	"4m": {Name: "4m", LowerFreq: 70000, UpperFreq: 70500},
	"2m": {Name: "2m", LowerFreq: 144000, UpperFreq: 146000},

	// UHF Band
	"70cm": {Name: "70cm", LowerFreq: 430000, UpperFreq: 440000},

	// Microwave Bands (in kHz for consistency)
	"23cm": {Name: "23cm", LowerFreq: 1240000, UpperFreq: 1325000},     // 1.24-1.325 GHz
	"13cm": {Name: "13cm", LowerFreq: 2300000, UpperFreq: 2450000},     // 2.3-2.45 GHz (simplified)
	"9cm":  {Name: "9cm", LowerFreq: 3400000, UpperFreq: 3410000},      // 3.4-3.41 GHz
	"6cm":  {Name: "6cm", LowerFreq: 5650000, UpperFreq: 5850000},      // 5.65-5.85 GHz (simplified)
	"3cm":  {Name: "3cm", LowerFreq: 10000000, UpperFreq: 10475000},    // 10-10.475 GHz (simplified)
}

// GetBandFromFrequency returns the band name for a given frequency in kHz
// Returns empty string if frequency is not in an amateur band
func GetBandFromFrequency(freqKHz float64) string {
	for name, band := range RSGBBandPlan {
		if freqKHz >= band.LowerFreq && freqKHz <= band.UpperFreq {
			return name
		}
	}
	return ""
}

// GetBandFromFrequencyMHz returns the band name for a given frequency in MHz
func GetBandFromFrequencyMHz(freqMHz float64) string {
	return GetBandFromFrequency(freqMHz * 1000)
}

// IsFrequencyInBand checks if a frequency (in kHz) falls within a specific band
func IsFrequencyInBand(freqKHz float64, bandName string) bool {
	band, exists := RSGBBandPlan[bandName]
	if !exists {
		return false
	}
	return freqKHz >= band.LowerFreq && freqKHz <= band.UpperFreq
}

// GetAllBandNames returns a slice of all band names
func GetAllBandNames() []string {
	names := make([]string, 0, len(RSGBBandPlan))
	for name := range RSGBBandPlan {
		names = append(names, name)
	}
	return names
}
