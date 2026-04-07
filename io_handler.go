package gobayes

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// NodeDefinition sert de structure intermédiaire pour le JSON
type NodeDefinition struct {
	Name    string    `json:"name"`
	States  []string  `json:"states"`
	Parents []string  `json:"parents"`
	CPD     []float64 `json:"cpd"`
}

type NetworkDefinition struct {
	Nodes []NodeDefinition `json:"nodes"`
}

// SaveToFile enregistre le réseau dans un fichier JSON
func (n *Network) SaveToFile(filename string) error {
	def := NetworkDefinition{}
	for _, node := range n.Nodes {
		nodeDef := NodeDefinition{
			Name:   node.Name,
			States: node.States,
			CPD:    node.CPD,
		}
		for _, p := range node.Parents {
			nodeDef.Parents = append(nodeDef.Parents, p.Name)
		}
		def.Nodes = append(def.Nodes, nodeDef)
	}

	data, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile charge un réseau à partir d'un fichier JSON
func LoadFromFile(filename string) (*Network, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var def NetworkDefinition
	if err := json.NewDecoder(file).Decode(&def); err != nil {
		return nil, err
	}

	net := NewNetwork()

	// 1. Créer tous les nœuds d'abord
	for _, nd := range def.Nodes {
		net.AddNode(nd.Name, nd.States)
	}

	// 2. Créer les liens et injecter les CPD
	for _, nd := range def.Nodes {
		for _, pName := range nd.Parents {
			net.AddEdge(pName, nd.Name)
		}
		node := net.Nodes[nd.Name]
		node.SetProbabilities(nd.CPD)
	}

	return net, nil
}