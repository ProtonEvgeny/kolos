package aig

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ProtonEvgeny/kolos/internal/model"
	"io"
	"strconv"
	"strings"
)

// parseHeader parses the header line of an AIG file format.
func parseHeader(line string) (M, I, L, O, A int, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 6 || parts[0] != "aig" {
		return 0, 0, 0, 0, 0, errors.New("invalid AIG header")
	}

	M, _ = strconv.Atoi(parts[1])
	I, _ = strconv.Atoi(parts[2])
	L, _ = strconv.Atoi(parts[3])
	O, _ = strconv.Atoi(parts[4])
	A, _ = strconv.Atoi(parts[5])

	if M != I+L+A {
		return 0, 0, 0, 0, 0, errors.New("invalid M value in AIG header")
	}

	return
}

// decodeDelta reads a delta-encoded unsigned integer from r.
func decodeDelta(r io.ByteReader) (uint64, error) {
	var x uint64
	var shift uint

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		x |= uint64(b&0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
	}

	return x, nil
}

// ParseAIG parses an AIG (And-Inverter Graph) from the provided reader.
// Returns a pointer to the constructed AIG structure or an error if parsing fails.
func ParseAIG(r io.Reader) (*model.AIG, error) {
	br := bufio.NewReader(r)

	header, _ := br.ReadString('\n')
	M, I, L, O, A, err := parseHeader(strings.TrimSpace(header))
	if err != nil {
		return nil, err
	}

	aig := &model.AIG{
		MaxVar:   M,
		Inputs:   make([]*model.Node, I),
		Outputs:  make([]*model.Node, 0, O),
		Latches:  make([]*model.Node, L),
		AndGates: make([]*model.Node, 0, A),
	}

	for i := 0; i < I; i++ {
		aig.Inputs[i] = &model.Node{
			ID:   i + 1,
			Type: model.Input,
		}
	}

	for i := 0; i < L; i++ {
		line, _, _ := br.ReadLine()
		nextStateLit, _ := strconv.Atoi(string(line))

		currentLit := 2 * (I + i + 1)
		currentVar := currentLit / 2

		latch := &model.Node{
			ID:       currentVar,
			Type:     model.Latch,
			Inverted: (currentLit & 1) == 1,
			NextState: &model.Node{
				ID:       nextStateLit / 2,
				Inverted: (nextStateLit & 1) == 1,
			},
		}

		aig.Latches[i] = latch
	}

	for i := 0; i < O; i++ {
		line, _, err := br.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("output %d: %w", i, err)
		}

		lit, err := strconv.Atoi(string(line))
		if err != nil {
			return nil, fmt.Errorf("invalid output literal %q: %w", line, err)
		}
		aig.Outputs = append(aig.Outputs, &model.Node{
			ID:       lit >> 1,
			Inverted: (lit & 1) == 1,
			Type:     model.Output,
		})
	}

	for i := 0; i < A; i++ {
		delta0, _ := decodeDelta(br)
		delta1, _ := decodeDelta(br)

		lhs := 2 * (I + L + 1 + i) // from AIG format specification
		rhs0 := lhs - int(delta0)
		rhs1 := rhs0 - int(delta1)

		aig.AndGates = append(aig.AndGates, &model.Node{
			ID:   (I + L + 1 + i),
			Type: model.AndGate,
			Children: []*model.Node{
				{ID: rhs0 >> 1, Inverted: (rhs0 & 1) == 1},
				{ID: rhs1 >> 1, Inverted: (rhs1 & 1) == 1},
			},
		})
	}

	return aig, nil
}

// LinkNodes resolves node IDs in aig.AndGates to pointers to nodes in the aig structure.
func LinkNodes(aig *model.AIG) {
	nodeMap := make(map[int]*model.Node)

	for _, node := range aig.Inputs {
		nodeMap[node.ID] = node
	}

	for _, node := range aig.Latches {
		nodeMap[node.ID] = node
	}

	for _, node := range aig.AndGates {
		nodeMap[node.ID] = node
	}

	for _, and := range aig.AndGates {
		for i, child := range and.Children {
			and.Children[i] = nodeMap[child.ID]
		}
	}

	for _, latch := range aig.Latches {
		if latch.NextState != nil {
			if resolved, exists := nodeMap[latch.NextState.ID]; exists {
				latch.NextState = resolved
			} else {
				latch.NextState = &model.Node{
					ID:       latch.NextState.ID,
					Inverted: latch.NextState.Inverted,
				}
			}
		}
	}
}
