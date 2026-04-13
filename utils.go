package gobayes

func union(a, b []string) []string {
	m := make(map[string]bool)
	var result []string

	// On ajoute les éléments de la première liste
	for _, item := range a {
		if !m[item] {
			m[item] = true
			result = append(result, item)
		}
	}

	// On ajoute les éléments de la deuxième liste s'ils n'existent pas encore
	for _, item := range b {
		if !m[item] {
			m[item] = true
			result = append(result, item)
		}
	}

	return result
}