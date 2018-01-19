package tachikoma

import "strings"

func Labels(source, base, counter string) []string {
	return []string{source, strings.ToUpper(base), strings.ToUpper(counter)}
}
