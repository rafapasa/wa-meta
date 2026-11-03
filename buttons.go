package wameta

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// Meta Cloud API limits for interactive reply buttons.
const (
	MaxReplyButtons     = 3
	MaxButtonTitleRunes = 20
)

// SendButtons sends an interactive WhatsApp message with reply buttons (1–3).
// Uses the same Graph /{phone-number-id}/messages endpoint as SendText.
func (c *Client) SendButtons(ctx context.Context, in SendButtonsInput) (*SendMessageResult, error) {
	if c.cfg.PhoneNumberID == "" || c.cfg.AccessToken == "" {
		return nil, fmt.Errorf("wameta: PhoneNumberID e AccessToken são obrigatórios")
	}
	if strings.TrimSpace(in.To) == "" {
		return nil, fmt.Errorf("wameta: destinatário (to) vazio")
	}
	if strings.TrimSpace(in.Body) == "" {
		return nil, fmt.Errorf("wameta: body dos botões vazio")
	}
	n := len(in.Buttons)
	if n < 1 || n > MaxReplyButtons {
		return nil, fmt.Errorf("wameta: quantidade de botões inválida: %d (1–%d)", n, MaxReplyButtons)
	}

	replies := make([]map[string]any, 0, n)
	for i, b := range in.Buttons {
		id := strings.TrimSpace(b.ID)
		title := strings.TrimSpace(b.Title)
		if title == "" {
			return nil, fmt.Errorf("wameta: botão %d sem title", i)
		}
		if id == "" {
			id = ButtonIDFromTitle(title)
		}
		if id == "" {
			return nil, fmt.Errorf("wameta: botão %d sem id", i)
		}
		// Cloud API: title max 20 characters
		if runeLen(title) > MaxButtonTitleRunes {
			return nil, fmt.Errorf("wameta: title do botão %d excede %d caracteres", i, MaxButtonTitleRunes)
		}
		replies = append(replies, map[string]any{
			"type": "reply",
			"reply": map[string]string{
				"id":    id,
				"title": title,
			},
		})
	}

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                in.To,
		"type":              "interactive",
		"interactive": map[string]any{
			"type": "button",
			"body": map[string]string{"text": in.Body},
			"action": map[string]any{
				"buttons": replies,
			},
		},
	}
	return c.send(ctx, payload)
}

// ButtonIDFromTitle derives a stable button id from a visible label
// (lowercase, non-letters/digits → underscore). Empty if nothing remains.
func ButtonIDFromTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return ""
	}
	var b strings.Builder
	lastUnderscore := false
	for _, r := range strings.ToLower(title) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore && b.Len() > 0 {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(b.String(), "_")
}

func runeLen(s string) int {
	return len([]rune(s))
}
