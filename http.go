package wameta

import "net/http"

// HTTPDoer abstracts http.Client for tests and injection.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
