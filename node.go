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