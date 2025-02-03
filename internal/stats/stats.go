package stats

import (
	"github.com/ProtonEvgeny/kolos/internal/model"
)

type Stats struct {
	Inputs   int
	Outputs  int
	Latches  int
	AndGates int
	MaxLevel int
}

// Calculate computes and returns statistics about the given AIG (And-Inverter Graph).
// It calculates the number of inputs, outputs, latches, AND gates, and the maximum level
// of the nodes in the graph. The level of a node is the longest path from any input
// node to the node itself. The function returns a Stats struct containing these statistics.
func Calculate(aig *model.AIG) Stats {
	stats := Stats{
		Inputs:   len(aig.Inputs),
		Outputs:  len(aig.Outputs),
		Latches:  len(aig.Latches),
		AndGates: len(aig.AndGates),
	}

	nodeLevels := make(map[*model.Node]int)

	for _, node := range aig.Inputs {
		nodeLevels[node] = 0
	}

	for _, latch := range aig.Latches {
		nodeLevels[latch] = 0
	}

	var computeLevel func(*model.Node) int
	computeLevel = func(node *model.Node) int {
		if level, exists := nodeLevels[node]; exists {
			return level
		}

		if node.Type == model.Output {
			if len(node.Children) == 0 {
				return 0
			}
			return computeLevel(node.Children[0])
		}

		maxChildLevel := -1
		for _, child := range node.Children {
			if cl := computeLevel(child); cl > maxChildLevel {
				maxChildLevel = cl
			}
		}

		level := maxChildLevel + 1
		if node.Type == model.Input || node.Type == model.Latch {
			level = 0
		}

		nodeLevels[node] = level
		return level
	}

	maxLevel := 0
	for _, and := range aig.AndGates {
		if level := computeLevel(and); level > maxLevel {
			maxLevel = level
		}
	}
	for _, out := range aig.Outputs {
		if level := computeLevel(out); level > maxLevel {
			maxLevel = level
		}
	}

	stats.MaxLevel = maxLevel
	return stats
}
