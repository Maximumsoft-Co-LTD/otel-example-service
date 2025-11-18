package utils

import "strings"

func BearerTokenFrom(h string) (string, bool) {
	if h == "" {
		return "", false
	}
	const p = "bearer "
	hl := strings.ToLower(h)
	if strings.HasPrefix(hl, p) {
		return strings.TrimSpace(h[len(p):]), true
	}
	const P = "Bearer "
	if strings.HasPrefix(h, P) {
		return strings.TrimSpace(h[len(P):]), true
	}
	return "", false
}
