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

// AddEdge crée un lien de causalité entre un parent et un enfant
func (n *Network) AddEdge(parentName, childName string) error {
	parent, okP := n.Nodes[parentName]
	child, okC := n.Nodes[childName]
	if !okP || !okC {
		return errors.New("parent ou enfant introuvable")
	}

	child.Parents = append(child.Parents, parent)
	parent.Children = append(parent.Children, child)

	// Vérification de cycle (Simple DFS)
	if n.hasCycle() {
		// On annule le lien si un cycle est détecté
		child.Parents = child.Parents[:len(child.Parents)-1]
		parent.Children = parent.Children[:len(parent.Children)-1]
		return errors.New("cycle détecté : un réseau bayésien doit être acyclique")
	}
	return nil
}

// SetProbabilities permet de remplir la table CPD. 
// Les probabilités doivent être fournies dans l'ordre des combinaisons d'états.
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

// hasCycle vérifie si le graphe contient un cycle
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
		return true // Cycle détecté !
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

// ToFactor convertit les probabilités d'un nœud en un objet Factor calculable
func (n *Node) ToFactor() *Factor {
	// 1. Collecter les variables impliquées : le nœud lui-même + ses parents
	vars := []string{n.Name}
	dims := make(map[string]int)
	dims[n.Name] = len(n.States)

	for _, p := range n.Parents {
		vars = append(vars, p.Name)
		dims[p.Name] = len(p.States)
	}

	// 2. Créer le facteur avec une copie des probabilités (CPD)
	// On copie pour éviter de modifier accidentellement les données d'origine
	cpdCopy := make([]float64, len(n.CPD))
	copy(cpdCopy, n.CPD)

	return &Factor{
		Variables: vars,
		Dims:      dims,
		Values:    cpdCopy,
	}
}