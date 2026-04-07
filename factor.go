package gobayes

// Factor représente une table de probabilités multidimensionnelle
// Multiply multiplie deux facteurs entre eux
func (f *Factor) Multiply(other *Factor) *Factor {
	newVars := union(f.Variables, other.Variables)
	newDims := make(map[string]int)
	size := 1
	for _, v := range newVars {
		dim, exists := f.Dims[v]
		if !exists {
			dim = other.Dims[v]
		}
		newDims[v] = dim
		size *= dim
	}

	result := &Factor{
		Variables: newVars,
		Dims:      newDims,
		Values:    make([]float64, size),
	}

	resStrides := result.GetStrides()
	fStrides := f.GetStrides()
	oStrides := other.GetStrides()

	// On parcourt chaque case de la nouvelle table de probabilités
	for i := 0; i < size; i++ {
		// 1. On trouve la combinaison d'états pour cet index
		states := result.IndexToStates(i, resStrides)
		
		// 2. On trouve l'index correspondant dans les deux facteurs d'origine
		idxF := f.StatesToIndex(states, fStrides)
		idxO := other.StatesToIndex(states, oStrides)
		
		// 3. On multiplie les probabilités
		result.Values[i] = f.Values[idxF] * other.Values[idxO]
	}

	return result
}

// Marginalize élimine une variable en sommant ses probabilités
func (f *Factor) Marginalize(variable string) *Factor {
	newVars := []string{}
	for _, v := range f.Variables {
		if v != variable {
			newVars = append(newVars, v)
		}
	}

	// Si on a éliminé la seule variable, on retourne un facteur scalaire
	// ... (gestion des cas limites)

	newDims := make(map[string]int)
	size := 1
	for _, v := range newVars {
		newDims[v] = f.Dims[v]
		size *= newDims[v]
	}

	result := &Factor{
		Variables: newVars,
		Dims:      newDims,
		Values:    make([]float64, size),
	}

	resStrides := result.GetStrides()
	fStrides := f.GetStrides()

	for i := 0; i < len(f.Values); i++ {
		states := f.IndexToStates(i, fStrides)
		// On projette l'état sur le nouveau facteur (sans la variable éliminée)
		resIdx := result.StatesToIndex(states, resStrides)
		result.Values[resIdx] += f.Values[i]
	}

	return result
}

// GetStrides calcule les pas (strides) pour chaque variable.
// Cela permet de savoir de combien on "saute" dans le slice pour changer l'état d'une variable.
func (f *Factor) GetStrides() map[string]int {
	strides := make(map[string]int)
	currentStride := 1
	// On parcourt à l'envers pour respecter l'ordre de stockage classique
	for i := len(f.Variables) - 1; i >= 0; i-- {
		v := f.Variables[i]
		strides[v] = currentStride
		currentStride *= f.Dims[v]
	}
	return strides
}

// IndexToStates convertit un index de slice plat en une map d'états (0, 1, 2...)
func (f *Factor) IndexToStates(index int, strides map[string]int) map[string]int {
	states := make(map[string]int)
	for _, v := range f.Variables {
		stride := strides[v]
		states[v] = (index / stride) % f.Dims[v]
	}
	return states
}

// StatesToIndex fait l'inverse : convertit une configuration d'états en un index unique
func (f *Factor) StatesToIndex(states map[string]int, strides map[string]int) int {
	index := 0
	for _, v := range f.Variables {
		index += states[v] * strides[v]
	}
	return index
}

// Reduce crée un nouveau facteur en fixant une variable à un index d'état spécifique.
// C'est l'opération mathématique pour prendre en compte une "preuve".
func (f *Factor) Reduce(variable string, stateIndex int) *Factor {
	// 1. Vérifier si la variable existe dans ce facteur
	exists := false
	for _, v := range f.Variables {
		if v == variable {
			exists = true
			break
		}
	}

	// Si la variable n'est pas dans ce facteur, on retourne le facteur tel quel
	if !exists {
		return f
	}

	// 2. Définir les nouvelles variables (on retire la variable fixée)
	newVars := []string{}
	for _, v := range f.Variables {
		if v != variable {
			newVars = append(newVars, v)
		}
	}

	newDims := make(map[string]int)
	size := 1
	for _, v := range newVars {
		newDims[v] = f.Dims[v]
		size *= newDims[v]
	}
	if size == 0 { size = 1 }

	result := &Factor{
		Variables: newVars,
		Dims:      newDims,
		Values:    make([]float64, size),
	}

	resStrides := result.GetStrides()
	fStrides := f.GetStrides()

	// 3. Remplir le nouveau facteur uniquement avec les lignes valides
	for i := 0; i < size; i++ {
		states := result.IndexToStates(i, resStrides)
		// On rajoute la variable fixée pour trouver l'index dans le facteur d'origine
		states[variable] = stateIndex
		
		oldIdx := f.StatesToIndex(states, fStrides)
		result.Values[i] = f.Values[oldIdx]
	}

	return result
}

func (f *Factor) Normalize() {
	sum := 0.0
	for _, v := range f.Values {
		sum += v
	}
	for i := range f.Values {
		f.Values[i] /= sum
	}
}