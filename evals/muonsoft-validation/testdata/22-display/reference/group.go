package scenario

func GroupErrors(issues []Issue) []Group {
	var result []Group
	positions := make(map[string]int)
	for _, issue := range issues {
		index, ok := positions[issue.Path]
		if !ok {
			index = len(result)
			positions[issue.Path] = index
			result = append(result, Group{Path: issue.Path})
		}
		result[index].Messages = append(result[index].Messages, issue.Message)
	}
	return result
}
