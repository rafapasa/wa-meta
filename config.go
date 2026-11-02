package wameta

// Config holds credentials and HTTP options for Graph API calls.
// Webhook security (AppSecret / VerifyToken) is also used by the package-level helpers.
type Config struct {
	// AppSecret used for HMAC-SHA256 of X-Hub-Signature-256.
	AppSecret string
	// VerifyToken expected on GET hub.verify_token during subscription challenge.
	VerifyToken string
	// PhoneNumberID sender phone number id on Graph API.
	PhoneNumberID string
	// AccessToken Graph API bearer token.
	AccessToken string
	// BaseURL Graph API host (default https://graph.facebook.com).
	BaseURL string
	// APIVersion Graph API version (default v21.0).
	APIVersion string
	// HTTPClient optional; when nil, http.DefaultClient is used.
	HTTPClient HTTPDoer
}
