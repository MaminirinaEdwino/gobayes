package gobayes

// Rule définit le poids d'une option selon un contexte
type Rule struct {
    TargetState  string
    ParentStates map[string]string
    Weight       float64
}

// GenerateCPD construit la table automatiquement à partir de règles
func (n *Node) GenerateCPD(rules []Rule) {
    // 1. Calculer toutes les combinaisons possibles (le produit cartésien)
    // 2. Pour chaque combinaison, appliquer les poids des règles correspondantes
    // 3. Normaliser pour que chaque bloc de probabilités somme à 1
    // 4. Injecter dans n.CPD
}