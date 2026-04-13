package gobayes

import (
	"errors"
	"fmt"
)

func NewNetwork() *Network {
	return &Network{Nodes: make(map[string]*Node)}
}

// AddNode crée et ajoute un nœud au réseau
func (n *Network) AddNode(name string, states []string) (*Node, error) {
	if _, exists := n.Nodes[name]; exists {
		return nil, fmt.Errorf("le nœud %s existe déjà", name)
	}
	node := &Node{
		Name:   name,
		States: states,
	}
	n.Nodes[name] = node
	return node, nil
}

func (n *Network) AddEdge(parentName, childName string) error {
	parent, okP := n.Nodes[parentName]
	child, okC := n.Nodes[childName]
	if !okP || !okC {
		return errors.New("parent ou enfant introuvable")
	}

	child.Parents = append(child.Parents, parent)
	parent.Children = append(parent.Children, child)

	if n.hasCycle() {
		child.Parents = child.Parents[:len(child.Parents)-1]
		parent.Children = parent.Children[:len(parent.Children)-1]
		return errors.New("cycle détecté : un réseau bayésien doit être acyclique")
	}
	return nil
}

func (node *Node) SetProbabilities(probs []float64) error {
	expectedSize := len(node.States)
	for _, p := range node.Parents {
		expectedSize *= len(p.States)
	}

	if len(probs) != expectedSize {
		return fmt.Errorf("taille de CPD incorrecte : attendu %d, reçu %d", expectedSize, len(probs))
	}

	node.CPD = probs
	return nil
}

func (n *Network) hasCycle() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range n.Nodes {
		if n.isCyclic(name, visited, recStack) {
			return true
		}
	}
	return false
}

func (n *Network) isCyclic(name string, visited, recStack map[string]bool) bool {
	if recStack[name] {
		return true
	}
	if visited[name] {
		return false
	}

	visited[name] = true
	recStack[name] = true

	// On explore tous les enfants
	for _, child := range n.Nodes[name].Children {
		if n.isCyclic(child.Name, visited, recStack) {
			return true
		}
	}

	recStack[name] = false // On retire de la pile après exploration
	return false
}

func (n *Node) ToFactor() *Factor {
	vars := []string{n.Name}
	dims := make(map[string]int)
	dims[n.Name] = len(n.States)

	for _, p := range n.Parents {
		vars = append(vars, p.Name)
		dims[p.Name] = len(p.States)
	}

	cpdCopy := make([]float64, len(n.CPD))
	copy(cpdCopy, n.CPD)

	return &Factor{
		Variables: vars,
		Dims:      dims,
		Values:    cpdCopy,
	}
}