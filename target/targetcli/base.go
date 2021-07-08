package targetcli

import "strings"

const (
	targetcliBinary = "targetcli"
)

func CleanResultString(prefix string, raw string) string {
	if prefix != "" {
		raw = strings.TrimPrefix(raw, prefix)
	}

	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, ".")

	return raw
}
