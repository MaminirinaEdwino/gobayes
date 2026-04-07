package gobayes

// Factor représente une table de probabilités multidimensionnelle
// Multiply multiplie deux facteurs entre eux
func (f *Factor) Multiply(other *Factor) *Factor {
    // 1. Déterminer les variables du nouveau facteur (Union des deux)
    newVars := union(f.Variables, other.Variables)
    newDims := make(map[string]int)
    
    // 2. Calculer la taille totale de la nouvelle table
    size := 1
    for _, v := range newVars {
        dim := f.Dims[v]
        if dim == 0 { dim = other.Dims[v] }
        newDims[v] = dim
        size *= dim
    }

    result := &Factor{
        Variables: newVars,
        Dims:      newDims,
        Values:    make([]float64, size),
    }

    // 3. Remplir les valeurs (C'est ici qu'on utilise les goroutines si besoin)
    // Pour chaque combinaison dans le nouveau facteur, on va chercher 
    // les valeurs correspondantes dans f et other.
    // ... logique d'indexation ...
    
    return result
}

// Marginalize élimine une variable en sommant ses probabilités
func (f *Factor) Marginalize(variable string) *Factor {
    // 1. Créer les nouvelles variables (f.Variables moins la variable à éliminer)
    // 2. Calculer les nouvelles dimensions
    // 3. Sommer les valeurs correspondantes
    // ...
    return result
}