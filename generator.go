package gobayes

import (
	"fmt"
	"log"
)


func (n *Node) GenerateCPD(rules []ScoreRule) {
	// 1. Calculer la taille totale de la table (Produit des états de tous les parents)
	totalParentCombos := 1
	for _, p := range n.Parents {
		totalParentCombos *= len(p.States)
	}

	// La table CPD finale aura (combinaisons parents * nb états du nœud actuel)
	newCPD := make([]float64, 0, totalParentCombos*len(n.States))

	// 2. Pour chaque combinaison possible des parents
	for i := 0; i < totalParentCombos; i++ {
		// Récupérer la configuration des parents pour cet index (ex: {TempsReel: "Oui", Equipe: "Solo"})
		parentStates := n.getParentStatesForIndex(i)

		// Initialiser les scores de chaque état de ce nœud à 1.0 (neutre)
		rowScores := make([]float64, len(n.States))
		for s := range rowScores {
			rowScores[s] = 1.0
		}
		log.Printf("Génération CPD pour le nœud : %s", n.Name)
		for _, rule := range rules {
			log.Printf("Règle trouvée : vise %s, état %s", rule.TargetNode, rule.TargetState)
		}
		for _, rule := range rules {
			// SÉCURITÉ : On n'applique la règle que si elle vise CE nœud précis
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

		// 4. Normalisation (La somme de la ligne doit être 1.0)
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
