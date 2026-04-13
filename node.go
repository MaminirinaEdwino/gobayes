package gobayes

func (n *Node) getStateIndex(stateName string) int {
	for i, s := range n.States {
		if s == stateName {
			return i
		}
	}
	return 0 
}
func (n *Node) getParentStatesForIndex(index int) map[string]string {
    states := make(map[string]string)
    remainder := index
    
    for i := len(n.Parents) - 1; i >= 0; i-- {
        p := n.Parents[i]
        stateIdx := remainder % len(p.States)
        states[p.Name] = p.States[stateIdx]
        remainder /= len(p.States)
    }
    return states
}
func (n *Node) ruleMatches(rule ScoreRule, parentStates map[string]string) bool {
    for varName, requiredState := range rule.Conditions {
        actualState, ok := parentStates[varName]
        if !ok || actualState != requiredState {
            return false
        }
    }
    return true
}