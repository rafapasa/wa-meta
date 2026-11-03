package wameta

// WebhookEvent is the root payload Meta posts to the business webhook.
type WebhookEvent struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

type WebhookEntry struct {
	ID      string          `json:"id"`
	Changes []WebhookChange `json:"changes"`
}

type WebhookChange struct {
	Field string       `json:"field"`
	Value WebhookValue `json:"value"`
}

type WebhookValue struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         WebhookMetadata  `json:"metadata"`
	Contacts         []WebhookContact `json:"contacts,omitempty"`
	Messages         []WebhookMessage `json:"messages,omitempty"`
	Statuses         []WebhookStatus  `json:"statuses,omitempty"`
}

type WebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type WebhookContact struct {
	Profile WebhookProfile `json:"profile"`
	WaID    string         `json:"wa_id"`
}

type WebhookProfile struct {
	Name string `json:"name"`
}

type WebhookMessage struct {
	From      string        `json:"from"`
	ID        string        `json:"id"`
	Timestamp string        `json:"timestamp"`
	Type      string        `json:"type"`
	Text      *MessageText  `json:"text,omitempty"`
	Image     *MessageMedia `json:"image,omitempty"`
	Audio     *MessageMedia `json:"audio,omitempty"`
	Video     *MessageMedia `json:"video,omitempty"`
	Document  *MessageMedia `json:"document,omitempty"`
}

type MessageText struct {
	Body string `json:"body"`
}

type MessageMedia struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type,omitempty"`
	SHA256   string `json:"sha256,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

type WebhookStatus struct {
	ID           string              `json:"id"`
	Status       string              `json:"status"`
	Timestamp    string              `json:"timestamp"`
	RecipientID  string              `json:"recipient_id"`
	Conversation *StatusConversation `json:"conversation,omitempty"`
	Pricing      *StatusPricing      `json:"pricing,omitempty"`
}

type StatusConversation struct {
	ID     string                `json:"id"`
	Origin *ConversationOrigin   `json:"origin,omitempty"`
}

type ConversationOrigin struct {
	Type string `json:"type"`
}

type StatusPricing struct {
	Billable     bool   `json:"billable"`
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
}

// SendTextInput is the input for sending a text message via Graph API.
type SendTextInput struct {
	To   string
	Body string
}

// SendMediaInput is the input for sending media by link or id.
type SendMediaInput struct {
	To       string
	Type     string // image|audio|video|document
	Link     string
	ID       string
	Caption  string
	Filename string
}

// SendMessageResult is a simplified Graph send response.
type SendMessageResult struct {
	MessagingProduct string                `json:"messaging_product"`
	Contacts         []SendMessageContact  `json:"contacts"`
	Messages         []SendMessageResponse `json:"messages"`
}

type SendMessageContact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type SendMessageResponse struct {
	ID string `json:"id"`
}

// MediaURL is returned by the media metadata endpoint before download.
type MediaURL struct {
	MessagingProduct string `json:"messaging_product"`
	URL              string `json:"url"`
	MimeType         string `json:"mime_type"`
	SHA256           string `json:"sha256"`
	FileSize         int64  `json:"file_size"`
	ID               string `json:"id"`
}


// ReplyButton is a WhatsApp interactive reply button (type=reply).
type ReplyButton struct {
	ID    string // payload id returned on button_reply
	Title string // label shown to the user (max 20 chars on Cloud API)
}

// SendButtonsInput sends an interactive message with up to 3 reply buttons.
type SendButtonsInput struct {
	To      string
	Body    string
	Buttons []ReplyButton
}
