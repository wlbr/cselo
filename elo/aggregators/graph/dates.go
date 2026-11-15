package graph

func optionalDate(date *string) string {
	if date == nil {
		return ""
	}
	return *date
}
