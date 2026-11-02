package wameta

import (
	"fmt"
	"strings"
)

// ValidateChallenge validates Meta webhook subscription handshake.
// On success it returns the challenge string that the HTTP handler must write in the response body (status 200).
// mode must be "subscribe" and token must match verifyToken.
func ValidateChallenge(verifyToken, mode, token, challenge string) (echo string, err error) {
	if strings.TrimSpace(verifyToken) == "" {
		return "", fmt.Errorf("wameta: verify token vazio")
	}
	if mode != "subscribe" {
		return "", fmt.Errorf("wameta: mode inválido %q", mode)
	}
	if token != verifyToken {
		return "", fmt.Errorf("wameta: verify token não confere")
	}
	if strings.TrimSpace(challenge) == "" {
		return "", fmt.Errorf("wameta: challenge vazio")
	}
	return challenge, nil
}
