package gobayes

import (
	"math"
	"testing"
)

func TestFactor_Multiply(t *testing.T) {
	// Facteur 1 : P(A)
	// A: [0 (False), 1 (True)]
	f1 := &Factor{
		Variables: []string{"A"},
		Dims:      map[string]int{"A": 2},
		Values:    []float64{0.6, 0.4}, // P(A=0)=0.6, P(A=1)=0.4
	}

	// Facteur 2 : P(B|A)
	// B: [0, 1]
	// Table: A=0 [0.7, 0.3], A=1 [0.2, 0.8]
	f2 := &Factor{
		Variables: []string{"A", "B"},
		Dims:      map[string]int{"A": 2, "B": 2},
		Values:    []float64{0.7, 0.3, 0.2, 0.8},
	}

	// Résultat attendu : P(A, B) = P(A) * P(B|A)
	// Index 0 (A=0, B=0) : 0.6 * 0.7 = 0.42
	// Index 1 (A=0, B=1) : 0.6 * 0.3 = 0.18
	// Index 2 (A=1, B=0) : 0.4 * 0.2 = 0.08
	// Index 3 (A=1, B=1) : 0.4 * 0.8 = 0.32
	expected := []float64{0.42, 0.18, 0.08, 0.32}

	result := f1.Multiply(f2)

	const epsilon = 1e-9 

	for i, v := range result.Values {
		// On calcule la différence absolue
		if math.Abs(v-expected[i]) > epsilon {
			t.Errorf("Valeur incorrecte à l'index %d : attendu %f, reçu %f", i, expected[i], v)
		}
	}
	
}