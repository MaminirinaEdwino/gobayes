package gobayes
// getStateIndex retrouve l'index entier d'un état à partir de son nom (ex: "Go" -> 1)
func (n *Node) getStateIndex(stateName string) int {
	for i, s := range n.States {
		if s == stateName {
			return i
		}
	}
	return 0 // Valeur par défaut si non trouvé
}

// getParentStatesForIndex décompose un index de ligne en une map d'états de parents
// func (n *Node) getParentStatesForIndex(index int) map[string]string {
// 	states := make(map[string]string)
// 	tempIndex := index

// 	// On parcourt les parents à l'envers pour correspondre à la logique des strides
// 	for i := len(n.Parents) - 1; i >= 0; i-- {
// 		p := n.Parents[i]
// 		stateIdx := tempIndex % len(p.States)
// 		states[p.Name] = p.States[stateIdx]
// 		tempIndex /= len(p.States)
// 	}
// 	return states
// }

func (n *Node) getParentStatesForIndex(index int) map[string]string {
    states := make(map[string]string)
    remainder := index
    
    // On parcourt les parents à l'envers pour correspondre à la logique des Strides
    for i := len(n.Parents) - 1; i >= 0; i-- {
        p := n.Parents[i]
        stateIdx := remainder % len(p.States)
        states[p.Name] = p.States[stateIdx]
        remainder /= len(p.States)
    }
    return states
}

// ruleMatches vérifie si les conditions d'une règle s'appliquent à la ligne actuelle
// func (n *Node) ruleMatches(rule ScoreRule, currentParentStates map[string]string) bool {
// 	for varName, requiredState := range rule.Conditions {
// 		if currentParentStates[varName] != requiredState {
// 			return false
// 		}
// 	}
// 	return true
// }

func (n *Node) ruleMatches(rule ScoreRule, parentStates map[string]string) bool {
    // On parcourt les conditions de la règle (ex: "TempsReel": "Oui")
    for varName, requiredState := range rule.Conditions {
        actualState, ok := parentStates[varName]
        
        // Si le parent n'existe pas ou que l'état ne match pas (ex: "Non" != "Oui")
        if !ok || actualState != requiredState {
            return false
        }
    }
    // Si toutes les conditions sont remplies
    return true
}