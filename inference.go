package gobayes

import "fmt"

// Query calcule la distribution de probabilité d'une variable cible,
// compte tenu des preuves fournies (evidence).
func (n *Network) Query(target string, evidence map[string]int) *Factor {
	// 1. Initialiser les facteurs à partir des CPDs de chaque nœud
	// factors := []*Factor{}
	// for _, node := range n.Nodes {
	// 	factors = append(factors, node.ToFactor())
	// }

	// // 2. Réduire les facteurs avec les preuves
	// for variable, state := range evidence {
	// 	for i, f := range factors {
	// 		factors[i] = f.Reduce(variable, state)
	// 		fmt.Println("partie reduce", factors[i])
	// 	}
	// }

	// // 3. Multiplier tous les facteurs entre eux
	// // result := factors[0]
	// // for i := 1; i < len(factors); i++ {
	// // 	result = result.Multiply(factors[i])
	// // }
	// // 3. Multiplier tous les facteurs entre eux de manière sécurisée
    // var result *Factor
    // for _, f := range factors {
    //     if result == nil {
    //         result = f
    //     } else {
    //         result = result.Multiply(f)
    //     }
    // }

	factors := []*Factor{}
    for _, node := range n.Nodes {
        f := node.ToFactor()
        // Réduction immédiate pour chaque facteur
        for variable, state := range evidence {
            f = f.Reduce(variable, state)
			fmt.Println("boucle reduce", f)
        }
        factors = append(factors, f)
    }

    // 3. Multiplication Sécurisée
    var result *Factor
    for _, f := range factors {
        // On ignore les facteurs "morts" (ceux qui n'ont plus de variables et valent 0)
        if len(f.Variables) == 0 && (len(f.Values) == 0 || f.Values[0] == 0) {
            continue
        }
        
        if result == nil {
            result = f
        } else {
            result = result.Multiply(f)
        }
		fmt.Println("Multiply", result)
    }

    // Si après tout ça result est nil, on crée un facteur uniforme pour éviter le crash
    if result == nil {
        return n.Nodes[target].ToFactor()
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