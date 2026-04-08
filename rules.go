package gobayes

type ScoreRule struct {
    TargetNode  string            `json:"target_node"`  // <-- Indispensable
    TargetState string            `json:"target_state"`
    Conditions  map[string]string `json:"conditions"`
    Weight      float64           `json:"weight"`
}

