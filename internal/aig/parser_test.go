package aig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProtonEvgeny/kolos/internal/aig"
	"github.com/ProtonEvgeny/kolos/internal/model"
)

func TestParseAIGFiles(t *testing.T) {
	files, err := filepath.Glob("../../test/epfl_benchmark/*.aig")
	if err != nil {
		t.Fatalf("Error while searching for .aig files: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {

			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("Error opening file %s: %v", file, err)
			}
			defer f.Close()

			// Parse AIG
			graph, err := aig.ParseAIG(f)
			if err != nil {
				t.Fatalf("Error parsing file %s: %v", file, err)
			}

			// Validate basic structure
			validateBasicStructure(t, graph)

			// Link nodes
			aig.LinkNodes(graph)
			validateGraphConnections(t, graph)
		})
	}
}

// validateBasicStructure check if the basic structure of the AIG is valid
func validateBasicStructure(t *testing.T, graph *model.AIG) {
	// Check if slices are initialized
	if graph.Inputs == nil {
		t.Error("Inputs slice is nil")
	}
	if graph.Outputs == nil {
		t.Error("Outputs slice is nil")
	}
	if graph.AndGates == nil {
		t.Error("AndGates slice is nil")
	}

	// Check Inputs are valid
	for i, input := range graph.Inputs {
		if input.Type != model.Input {
			t.Errorf("Input node %d has wrong type: %v", i, input.Type)
		}
		if input.ID != i+1 {
			t.Errorf("Input node %d has wrong ID: %d", i, input.ID)
		}
	}

	// Check Latches are valid
	for i, latch := range graph.Latches {
		if latch.Type != model.Latch {
			t.Errorf("Latch node %d has wrong type: %v", i, latch.Type)
		}
		if latch.NextState == nil {
			t.Errorf("Latch node %d has no next state", i)
		}
	}

	// Check Outputs are valid
	for i, output := range graph.Outputs {
		if output.Type != model.Output {
			t.Errorf("Output node %d has wrong type: %v", i, output.Type)
		}
	}

	// Check AND-gates
	for i, and := range graph.AndGates {
		if and.Type != model.AndGate {
			t.Errorf("AND gate %d has wrong type: %v", i, and.Type)
		}
		if len(and.Children) != 2 {
			t.Errorf("AND gate %d has wrong number of children: %d", i, len(and.Children))
		}
	}
}

// validateGraphConnections checks the correctness of the connections in the graph
func validateGraphConnections(t *testing.T, graph *model.AIG) {
	// Create a map to store all nodes for quick lookup
	nodeMap := make(map[int]*model.Node)

	// Add inputs to the map
	for _, input := range graph.Inputs {
		nodeMap[input.ID] = input
	}

	// Add AND-gates to the map
	for _, and := range graph.AndGates {
		nodeMap[and.ID] = and
	}

	// Check the connections of AND-gates
	for _, and := range graph.AndGates {
		for i, child := range and.Children {
			if child == nil {
				t.Errorf("AND gate %d has nil child at position %d", and.ID, i)
				continue
			}

			// Check that the child node exists in the map
			if _, exists := nodeMap[child.ID]; !exists {
				t.Errorf("AND gate %d references non-existent node ID %d", and.ID, child.ID)
			}
		}
	}

	// Check the connections of outputs
	for i, output := range graph.Outputs {
		if output.ID > 0 {
			if _, exists := nodeMap[output.ID]; !exists {
				t.Errorf("Output %d references non-existent node ID %d", i, output.ID)
			}
		}
	}

	// Check the connections of latches
	for i, latch := range graph.Latches {
		if latch.NextState == nil {
			t.Errorf("Latch %d has nil NextState", i)
			continue
		}

		// Check that the NextState node exists in the map
		if _, exists := nodeMap[latch.NextState.ID]; !exists {
			t.Errorf("Latch %d references non-existent NextState ID %d", i, latch.NextState.ID)
		}
	}
}
