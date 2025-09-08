package commands

import "strings"

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		".", "\\.",
		"-", "\\-",
		"=", "\\=",
		"(", "\\(",
		")", "\\)",
	)
	return replacer.Replace(text)
}
