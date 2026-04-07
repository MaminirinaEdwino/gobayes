package gobayes

// ScoreRule définit le poids d'un état selon une condition
type ScoreRule struct {
	TargetState string            // ex: "Go"
	Conditions  map[string]string // ex: {"TempsReel": "Oui"}
	Weight      float64           // ex: 10.0 (un poids élevé = probabilité haute)
}