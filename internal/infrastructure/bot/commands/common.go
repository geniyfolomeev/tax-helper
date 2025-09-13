package commands

import "strings"

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		".", "\\.",
		"-", "\\-",
		"=", "\\=",
		"(", "\\(",
		")", "\\)",
		"_", "\\_",
	)
	return replacer.Replace(text)
}
