package wameta

import (
	"encoding/json"
	"fmt"
)

// ParseWebhook decodes Meta webhook JSON. Pure function: no HTTP, no ACK.
func ParseWebhook(body []byte) (*WebhookEvent, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("wameta: body vazio")
	}
	var ev WebhookEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("wameta: parse webhook: %w", err)
	}
	if ev.Object == "" && len(ev.Entry) == 0 {
		return nil, fmt.Errorf("wameta: payload webhook vazio ou inválido")
	}
	return &ev, nil
}

// ParseWebhookVerified validates signature then parses. Still pure regarding HTTP status codes:
// the host app decides when to respond 200.
func ParseWebhookVerified(appSecret, signatureHeader string, body []byte) (*WebhookEvent, error) {
	if err := ValidateSignature(appSecret, signatureHeader, body); err != nil {
		return nil, err
	}
	return ParseWebhook(body)
}
