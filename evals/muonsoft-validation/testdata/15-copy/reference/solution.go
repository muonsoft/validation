package scenario

func CopyMessages(messages map[string][]string) map[string][]string {
	if messages == nil {
		return nil
	}
	result := make(map[string][]string, len(messages))
	for field, messages := range messages {
		if messages == nil {
			result[field] = nil
			continue
		}
		result[field] = append(make([]string, 0, len(messages)), messages...)
	}
	return result
}
