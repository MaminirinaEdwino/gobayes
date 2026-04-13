package gobayes

import (
	"fmt"
	"log"
)


func (n *Node) GenerateCPD(rules []ScoreRule) {
	totalParentCombos := 1
	for _, p := range n.Parents {
		totalParentCombos *= len(p.States)
	}
	newCPD := make([]float64, 0, totalParentCombos*len(n.States))
	for i := 0; i < totalParentCombos; i++ {
		parentStates := n.getParentStatesForIndex(i)

		rowScores := make([]float64, len(n.States))
		for s := range rowScores {
			rowScores[s] = 1.0
		}
		log.Printf("Génération CPD pour le nœud : %s", n.Name)
		for _, rule := range rules {
			log.Printf("Règle trouvée : vise %s, état %s", rule.TargetNode, rule.TargetState)
		}
		for _, rule := range rules {
			fmt.Println(n)
			if rule.TargetNode != n.Name {
				continue
			}

			if n.ruleMatches(rule, parentStates) {
				stateIdx := n.getStateIndex(rule.TargetState)
				if stateIdx != -1 {
					rowScores[stateIdx] *= rule.Weight
				}
			}
		}
		sum := 0.0
		for _, s := range rowScores {
			sum += s
		}
		for _, s := range rowScores {
			newCPD = append(newCPD, s/sum)
		}
	}
	fmt.Println("new cpd", newCPD)
	n.CPD = newCPD
}
