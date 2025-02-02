package stats

import (
	"github.com/ProtonEvgeny/kolos/internal/model"
)

type Stats struct {
	Inputs            int
	Outputs           int
	Latches           int
	AndGates          int
	MaxLevel          int
	LevelDistribution map[int]int
}

func Calculate(aig *model.AIG) Stats {
	stats := Stats{
		Inputs:            len(aig.Inputs),
		Outputs:           len(aig.Outputs),
		Latches:           len(aig.Latches),
		AndGates:          len(aig.AndGates),
		LevelDistribution: make(map[int]int),
	}

	nodeLevels := make(map[*model.Node]int)

	for _, node := range aig.Inputs {
		nodeLevels[node] = 0
		stats.LevelDistribution[0]++
	}

	for _, latch := range aig.Latches {
		nodeLevels[latch] = 0
		stats.LevelDistribution[0]++
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

		if node.Type == model.AndGate {
			stats.LevelDistribution[level]++
		}

		return level
	}

	for _, and := range aig.AndGates {
		level := computeLevel(and)
		if level > stats.MaxLevel {
			stats.MaxLevel = level
		}
	}
	for _, out := range aig.Outputs {
		level := computeLevel(out)
		if level > stats.MaxLevel {
			stats.MaxLevel = level
		}
	}

	return stats
}
