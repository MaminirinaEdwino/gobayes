package gobayes

// Query calcule la distribution de probabilité d'une variable cible, 
// compte tenu des preuves fournies (evidence).
func (n *Network) Query(target string, evidence map[string]int) *Factor {
	// 1. Initialiser les facteurs à partir des CPDs de chaque nœud
	factors := []*Factor{}
	for _, node := range n.Nodes {
		factors = append(factors, node.ToFactor())
	}

	// 2. Réduire les facteurs avec les preuves
	for variable, state := range evidence {
		for i, f := range factors {
			factors[i] = f.Reduce(variable, state)
		}
	}

	// 3. Multiplier tous les facteurs entre eux
	result := factors[0]
	for i := 1; i < len(factors); i++ {
		result = result.Multiply(factors[i])
	}

	// 4. Marginaliser toutes les variables sauf la cible
	for _, v := range result.Variables {
		if v != target {
			result = result.Marginalize(v)
		}
	}

	// 5. Normaliser le résultat (la somme doit valoir 1)
	result.Normalize()

	return result
}