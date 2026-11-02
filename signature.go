package wameta

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ValidateSignature checks X-Hub-Signature-256 (sha256=<hex>).
// Returns nil when the signature matches appSecret over body.
func ValidateSignature(appSecret, signatureHeader string, body []byte) error {
	if strings.TrimSpace(appSecret) == "" {
		return fmt.Errorf("wameta: app secret vazio")
	}
	sig := strings.TrimSpace(signatureHeader)
	if sig == "" {
		return fmt.Errorf("wameta: signature header vazio")
	}
	const prefix = "sha256="
	if !strings.HasPrefix(strings.ToLower(sig), prefix) {
		// accept both sha256= and SHA256=
		if !strings.HasPrefix(sig, "sha256=") && !strings.HasPrefix(sig, "SHA256=") {
			return fmt.Errorf("wameta: signature header inválido")
		}
	}
	// normalize
	i := strings.Index(sig, "=")
	if i < 0 || i+1 >= len(sig) {
		return fmt.Errorf("wameta: signature header inválido")
	}
	gotHex := strings.TrimSpace(sig[i+1:])
	mac := hmac.New(sha256.New, []byte(appSecret))
	_, _ = mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(strings.ToLower(gotHex)), []byte(strings.ToLower(expected))) {
		return fmt.Errorf("wameta: assinatura inválida")
	}
	return nil
}
