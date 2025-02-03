package stats_test

import (
	"github.com/ProtonEvgeny/kolos/internal/aig"
	"github.com/ProtonEvgeny/kolos/internal/stats"
	"os"
	"path/filepath"
	"testing"
)

type ExpectedStats struct {
	MaxLevel int
	Inputs   int
	Outputs  int
	Latches  int
	AndGates int
}

func TestStats(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		expected ExpectedStats
	}{
		{
			name: "adder",
			file: "adder.aig",
			expected: ExpectedStats{
				MaxLevel: 255,
				Inputs:   256,
				Outputs:  129,
				Latches:  0,
				AndGates: 1020,
			},
		},
		{
			name: "multiplier",
			file: "multiplier.aig",
			expected: ExpectedStats{
				MaxLevel: 274,
				Inputs:   128,
				Outputs:  128,
				Latches:  0,
				AndGates: 27062,
			},
		},
		{
			name: "mem_ctrl",
			file: "mem_ctrl.aig",
			expected: ExpectedStats{
				MaxLevel: 114,
				Inputs:   1204,
				Outputs:  1231,
				Latches:  0,
				AndGates: 46836,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join("../../test/epfl_benchmark", tc.file)
			file, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("Error opening file %s: %v", tc.file, err)
			}
			defer file.Close()

			graph, err := aig.ParseAIG(file)
			if err != nil {
				t.Fatalf("Error parsing file %s: %v", tc.file, err)
			}
			aig.LinkNodes(graph)

			actual := stats.Calculate(graph)

			if actual.MaxLevel != tc.expected.MaxLevel {
				t.Errorf("MaxLevel: expected %d, got %d",
					tc.expected.MaxLevel, actual.MaxLevel)
			}
			if actual.Inputs != tc.expected.Inputs {
				t.Errorf("Inputs: expected %d, got %d",
					tc.expected.Inputs, actual.Inputs)
			}
			if actual.Outputs != tc.expected.Outputs {
				t.Errorf("Outputs: expected %d, got %d",
					tc.expected.Outputs, actual.Outputs)
			}
			if actual.Latches != tc.expected.Latches {
				t.Errorf("Latches: expected %d, got %d",
					tc.expected.Latches, actual.Latches)
			}
			if actual.AndGates != tc.expected.AndGates {
				t.Errorf("AndGates: expected %d, got %d",
					tc.expected.AndGates, actual.AndGates)
			}
		})
	}
}
