package lightrouter

func formatHeaders(headers map[string]string) string {
	result := ""
	for key, value := range headers {
		result += key + ": " + value + "\r\n"
	}
	return result
}
