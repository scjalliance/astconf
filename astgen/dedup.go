package astgen

func dedupStringSlice(values []string) []string {
	if len(values) == 0 {
		return values
	}

	seen := make(map[string]struct{}, len(values))
	output := make([]string, 0, len(values))

	for _, value := range values {
		if _, present := seen[value]; !present {
			output = append(output, value)
			seen[value] = struct{}{}
		}
	}

	return output
}
