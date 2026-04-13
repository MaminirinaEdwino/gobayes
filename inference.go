package gobayes

import "fmt"

func (n *Network) Query(target string, evidence map[string]int) *Factor {
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

    if result == nil {
        return n.Nodes[target].ToFactor()
    }

	for _, v := range result.Variables {
		if v != target {
			result = result.Marginalize(v)
		}
	}

	result.Normalize()
	return result
}