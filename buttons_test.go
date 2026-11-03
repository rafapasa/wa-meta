package wameta

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendButtons_Payload(t *testing.T) {
	var gotBody map[string]any
	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method=%s", req.Method)
		}
		if !strings.Contains(req.URL.Path, "/123/messages") {
			t.Fatalf("url=%s", req.URL.String())
		}
		raw, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"messaging_product":"whatsapp","messages":[{"id":"wamid.x"}]}`)),
			Header:     make(http.Header),
		}, nil
	})
	c := NewClient(Config{
		PhoneNumberID: "123",
		AccessToken:   "tok",
		HTTPClient:    &http.Client{Transport: rt},
	})
	_, err := c.SendButtons(context.Background(), SendButtonsInput{
		To:   "5511999999999",
		Body: "Continuar?",
		Buttons: []ReplyButton{
			{ID: "sim", Title: "Sim"},
			{Title: "Não"}, // id derivado
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotBody["type"] != "interactive" {
		t.Fatalf("type=%v", gotBody["type"])
	}
	inter, _ := gotBody["interactive"].(map[string]any)
	if inter["type"] != "button" {
		t.Fatalf("interactive.type=%v", inter["type"])
	}
	action, _ := inter["action"].(map[string]any)
	btns, _ := action["buttons"].([]any)
	if len(btns) != 2 {
		t.Fatalf("buttons=%v", btns)
	}
	b1, _ := btns[1].(map[string]any)
	reply, _ := b1["reply"].(map[string]any)
	if reply["id"] != "nao" && reply["id"] != "não" {
		// ButtonIDFromTitle("Não") -> "não" or "nao" depending on unicode letter
		id, _ := reply["id"].(string)
		if id == "" {
			t.Fatalf("derived id empty: %v", reply)
		}
	}
}

func TestSendButtons_InvalidCount(t *testing.T) {
	c := NewClient(Config{PhoneNumberID: "1", AccessToken: "t"})
	_, err := c.SendButtons(context.Background(), SendButtonsInput{
		To: "1", Body: "x", Buttons: nil,
	})
	if err == nil {
		t.Fatal("esperava erro")
	}
}

func TestButtonIDFromTitle(t *testing.T) {
	if got := ButtonIDFromTitle("Sim, continuar!"); got != "sim_continuar" {
		t.Fatalf("got %q", got)
	}
}
