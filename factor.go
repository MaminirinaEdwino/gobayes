package gobayes

import "math"

func (f *Factor) Multiply(other *Factor) *Factor {
	if len(f.Values) == 0 { return other }
    if len(other.Values) == 0 { return f }
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

	if size == 0 { return result }

	resStrides := result.GetStrides()
	fStrides := f.GetStrides()
	oStrides := other.GetStrides()

	for i := 0; i < size; i++ {
		states := result.IndexToStates(i, resStrides)
		
		idxF := f.StatesToIndex(states, fStrides)
		idxO := other.StatesToIndex(states, oStrides)
		
		result.Values[i] = f.Values[idxF] * other.Values[idxO]
	}

	return result
}

func (f *Factor) Marginalize(variable string) *Factor {
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

func (f *Factor) GetStrides() map[string]int {
	strides := make(map[string]int)
	if len(f.Variables) == 0 {
        return strides
    }
	currentStride := 1
	for i := len(f.Variables) - 1; i >= 0; i-- {
		v := f.Variables[i]
		strides[v] = currentStride
		currentStride *= f.Dims[v]
	}
	return strides
}

func (f *Factor) IndexToStates(index int, strides map[string]int) map[string]int {
	states := make(map[string]int)
	if len(f.Variables) == 0 {
        return states 
    }
	for _, v := range f.Variables {
		stride := strides[v]
		states[v] = (index / stride) % f.Dims[v]
	}
	return states
}

func (f *Factor) StatesToIndex(states map[string]int, strides map[string]int) int {
	index := 0
	for _, v := range f.Variables {
		index += states[v] * strides[v]
	}
	return index
}
func (f *Factor) Reduce(variable string, stateIndex int) *Factor {
    exists := false
    for _, v := range f.Variables {
        if v == variable {
            exists = true
            break
        }
    }
    if !exists {
        return f
    }

    newVars := []string{}
    newDims := make(map[string]int)
    size := 1
    for _, v := range f.Variables {
        if v != variable {
            newVars = append(newVars, v)
            newDims[v] = f.Dims[v]
            size *= f.Dims[v]
        }
    }

    if size <= 0 {
        size = 1
    }

    result := &Factor{
        Variables: newVars,
        Dims:      newDims,
        Values:    make([]float64, size),
    }

    fStrides := f.GetStrides()

    if len(newVars) == 0 {
        tempStates := map[string]int{variable: stateIndex}
        oldIdx := f.StatesToIndex(tempStates, fStrides)
        
        if oldIdx < len(f.Values) {
            result.Values[0] = f.Values[oldIdx]
        }
    } else {
        resStrides := result.GetStrides()
        for i := 0; i < size; i++ {
            states := result.IndexToStates(i, resStrides)
            states[variable] = stateIndex
            
            oldIdx := f.StatesToIndex(states, fStrides)
            if oldIdx < len(f.Values) && i < len(result.Values) {
                result.Values[i] = f.Values[oldIdx]
            }
        }
    }

    return result
}

func (f *Factor) Normalize() {
    sum := 0.0
    for _, v := range f.Values {
        sum += v
    }

    if sum == 0 || math.IsNaN(sum) {
        if len(f.Values) > 0 {
            uniform := 1.0 / float64(len(f.Values))
            for i := range f.Values {
                f.Values[i] = uniform
            }
        }
        return
    }

    for i := range f.Values {
        f.Values[i] /= sum
    }
}