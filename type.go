package gobayes

// StateIndex représente l'index d'un état (0, 1, 2...)
type StateIndex int

// Node représente un sommet du graphe
type Node struct {
	Name     string
	States   []string
	Parents  []*Node
	Children []*Node
	// CPD stockée sous forme de slice plat pour la performance
	CPD      []float64 
}

// Network gère l'ensemble des nœuds
type Network struct {
	Nodes map[string]*Node
}