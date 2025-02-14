package model

type NodeType int

const (
	Input NodeType = iota
	Output
	AndGate
	Latch
)

type Node struct {
	ID           int
	Type         NodeType
	Inverted     bool
	Children     []*Node
	NextState    *Node
	InitialState bool // add for simulation support (now useless)
}

type AIG struct {
	Inputs   []*Node
	Outputs  []*Node
	Latches  []*Node
	AndGates []*Node
	MaxVar   int
}
