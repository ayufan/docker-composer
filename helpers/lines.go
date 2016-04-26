package helpers

import "strings"

func Lines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	text := string(data)
	text = strings.TrimSpace(text)
	return strings.Split(text, "\n")
}
