package gobayes

// ScoreRule définit le poids d'un état selon une condition
// type ScoreRule struct {
//     TargetState string           `json:"target_state"`
//     Conditions  map[string]string `json:"conditions"`
//     Weight      float64          `json:"weight"`
// }
type ScoreRule struct {
    TargetNode  string            `json:"target_node"`  // <-- Indispensable
    TargetState string            `json:"target_state"`
    Conditions  map[string]string `json:"conditions"`
    Weight      float64           `json:"weight"`
}

// GenerateAutomatedCPD calcule la table de probabilités à partir de règles de scoring
func (n *Node) GenerateAutomatedCPD(rules []ScoreRule) {
	// 1. Calculer la taille totale de la table
	totalRows := 1
	for _, p := range n.Parents {
		totalRows *= len(p.States)
	}

	newCPD := make([]float64, 0)

	// 2. Pour chaque combinaison de parents
	for i := 0; i < totalRows; i++ {
		// Récupérer la configuration actuelle des parents pour cette ligne
		parentStates := n.getParentStatesForIndex(i)
		
		rowScores := make([]float64, len(n.States))
		// Initialiser avec un score de base (pour éviter les probabilités à 0)
		for s := range rowScores { rowScores[s] = 1.0 }

		// 3. Appliquer les règles
		for _, rule := range rules {
			if n.ruleMatches(rule, parentStates) {
				stateIdx := n.getStateIndex(rule.TargetState)
				rowScores[stateIdx] *= rule.Weight
			}
		}

		// 4. Normaliser la ligne (Somme = 1.0)
		sum := 0.0
		for _, s := range rowScores { sum += s }
		for _, s := range rowScores {
			newCPD = append(newCPD, s/sum)
		}
	}
	n.CPD = newCPD
}