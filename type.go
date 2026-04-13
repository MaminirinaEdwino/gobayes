package gobayes

type StateIndex int

type Node struct {
	Name     string
	States   []string
	Parents  []*Node
	Children []*Node
	CPD      []float64 
}

type Network struct {
	Nodes map[string]*Node
}

type Factor struct {
	Variables []string           
	Dims      map[string]int     
	Values    []float64          
}