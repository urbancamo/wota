package main

import (
	"fmt"
	"sort"

	"github.com/urbancamo/wota/sotautils"
)

func main() {
	fmt.Println("RSGB Amateur Radio Band Plan")
	fmt.Println("============================")

	// Get all band names and sort them by frequency
	type bandInfo struct {
		name      string
		lowerFreq float64
	}

	var bands []bandInfo
	for name, band := range sotautils.RSGBBandPlan {
		bands = append(bands, bandInfo{name, band.LowerFreq})
	}

	// Sort by lower frequency
	sort.Slice(bands, func(i, j int) bool {
		return bands[i].lowerFreq < bands[j].lowerFreq
	})

	// Display all bands
	for _, b := range bands {
		band := sotautils.RSGBBandPlan[b.name]
		fmt.Printf("%-6s  %10.1f - %10.1f kHz", band.Name, band.LowerFreq, band.UpperFreq)

		// Show in MHz for easier reading of VHF/UHF/Microwave
		if band.LowerFreq >= 50000 {
			fmt.Printf("  (%6.3f - %7.3f MHz)", band.LowerFreq/1000, band.UpperFreq/1000)
		}
		fmt.Println()
	}

	fmt.Println("\n\nExample Frequency Lookups:")
	fmt.Println("==========================")

	// Test some example frequencies
	testFrequencies := []struct {
		freq float64
		unit string
	}{
		{7.100, "MHz"},
		{14.200, "MHz"},
		{145.500, "MHz"},
		{433.500, "MHz"},
		{1296.100, "MHz"},
	}

	for _, test := range testFrequencies {
		band := sotautils.GetBandFromFrequencyMHz(test.freq)
		if band != "" {
			fmt.Printf("%8.3f %-3s is in the %-6s band\n", test.freq, test.unit, band)
		} else {
			fmt.Printf("%8.3f %-3s is not in an amateur band\n", test.freq, test.unit)
		}
	}

	fmt.Println("\n\nCheck Specific Bands:")
	fmt.Println("=====================")

	// Check if specific frequencies are in specific bands
	checks := []struct {
		freq float64
		band string
	}{
		{14.200, "20m"},
		{7.100, "40m"},
		{7.100, "20m"}, // This should be false
		{145.500, "2m"},
	}

	for _, check := range checks {
		inBand := sotautils.IsFrequencyInBand(check.freq*1000, check.band)
		result := "NO"
		if inBand {
			result = "YES"
		}
		fmt.Printf("Is %7.3f MHz in %-6s? %s\n", check.freq, check.band, result)
	}
}
