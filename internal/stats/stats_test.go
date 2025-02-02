package stats_test

import (
	"testing"
	"github.com/ProtonEvgeny/kolos/internal/model"
	"github.com/ProtonEvgeny/kolos/internal/stats"
)

func TestStats(t *testing.T) {
	// AND(x1, x2) -> output
	andNode := &model.Node{
        ID: 3, 
        Type: model.AndGate,
        Children: []*model.Node{
            {ID: 1, Type: model.Input}, // x1
            {ID: 2, Type: model.Input}, // x2
        },
    }

    aig := &model.AIG{
        Inputs: []*model.Node{
            {ID: 1, Type: model.Input},
            {ID: 2, Type: model.Input},
        },
        AndGates: []*model.Node{andNode},
        Outputs:  []*model.Node{
			{ID: 3, Type: model.Output, Children: []*model.Node{andNode}},
		},
    }

	stats := stats.Calculate(aig)

	if stats.MaxLevel != 1 {
		t.Errorf("Expected max level 1, got %d", stats.MaxLevel)
	}
	if stats.LevelDistribution[0] != 2 {
		t.Errorf("Expected 2 nodes at level 0, got %d", stats.LevelDistribution[0])
	}
	if stats.LevelDistribution[1] != 1 {
		t.Errorf("Expected 1 nodes at level 1, got %d", stats.LevelDistribution[1])
	}
}