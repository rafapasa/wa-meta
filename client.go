package wameta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client performs authenticated Graph API operations (send / media).
// Signature, challenge and parse are package-level helpers and do not require Client.
type Client struct {
	cfg  Config
	http HTTPDoer
}

// NewClient returns a Graph API client. HTTPClient may be nil (uses http.DefaultClient).
func NewClient(cfg Config) *Client {
	h := cfg.HTTPClient
	if h == nil {
		h = http.DefaultClient
	}
	return &Client{cfg: cfg, http: h}
}

func (c *Client) base() string {
	base := strings.TrimRight(c.cfg.BaseURL, "/")
	if base == "" {
		base = "https://graph.facebook.com"
	}
	ver := strings.Trim(c.cfg.APIVersion, "/")
	if ver == "" {
		ver = "v21.0"
	}
	return base + "/" + ver
}

// SendText sends a WhatsApp text message via Cloud API.
func (c *Client) SendText(ctx context.Context, in SendTextInput) (*SendMessageResult, error) {
	if c.cfg.PhoneNumberID == "" || c.cfg.AccessToken == "" {
		return nil, fmt.Errorf("wameta: PhoneNumberID e AccessToken são obrigatórios")
	}
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                in.To,
		"type":              "text",
		"text":              map[string]string{"body": in.Body},
	}
	return c.send(ctx, payload)
}

// SendMedia sends image|audio|video|document by link or media id.
func (c *Client) SendMedia(ctx context.Context, in SendMediaInput) (*SendMessageResult, error) {
	if c.cfg.PhoneNumberID == "" || c.cfg.AccessToken == "" {
		return nil, fmt.Errorf("wameta: PhoneNumberID e AccessToken são obrigatórios")
	}
	mediaType := strings.ToLower(strings.TrimSpace(in.Type))
	if mediaType == "" {
		return nil, fmt.Errorf("wameta: tipo de mídia vazio")
	}
	mediaObj := map[string]any{}
	if in.Link != "" {
		mediaObj["link"] = in.Link
	}
	if in.ID != "" {
		mediaObj["id"] = in.ID
	}
	if in.Caption != "" {
		mediaObj["caption"] = in.Caption
	}
	if in.Filename != "" {
		mediaObj["filename"] = in.Filename
	}
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                in.To,
		"type":              mediaType,
		mediaType:           mediaObj,
	}
	return c.send(ctx, payload)
}

func (c *Client) send(ctx context.Context, payload map[string]any) (*SendMessageResult, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/messages", c.base(), c.cfg.PhoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("wameta: send %d: %s", resp.StatusCode, truncate(string(body), 400))
	}
	var out SendMessageResult
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("wameta: decode send: %w", err)
	}
	return &out, nil
}

// GetMediaURL resolves media id to a temporary download URL and metadata.
func (c *Client) GetMediaURL(ctx context.Context, mediaID string) (*MediaURL, error) {
	if c.cfg.AccessToken == "" {
		return nil, fmt.Errorf("wameta: AccessToken obrigatório")
	}
	if mediaID == "" {
		return nil, fmt.Errorf("wameta: mediaID vazio")
	}
	url := fmt.Sprintf("%s/%s", c.base(), mediaID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("wameta: media meta %d: %s", resp.StatusCode, truncate(string(body), 400))
	}
	var out MediaURL
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadMedia downloads full media bytes from a URL returned by GetMediaURL.
func (c *Client) DownloadMedia(ctx context.Context, mediaURL string) ([]byte, error) {
	if mediaURL == "" {
		return nil, fmt.Errorf("wameta: mediaURL vazio")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaURL, nil)
	if err != nil {
		return nil, err
	}
	if c.cfg.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("wameta: download %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	return body, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
