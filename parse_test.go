package wameta

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWebhook_TextFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "webhook_text.json"))
	if err != nil {
		t.Fatal(err)
	}
	ev, err := ParseWebhook(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Object != "whatsapp_business_account" {
		t.Fatalf("object=%s", ev.Object)
	}
	if len(ev.Entry) == 0 || len(ev.Entry[0].Changes) == 0 {
		t.Fatal("entry/changes vazios")
	}
	msgs := ev.Entry[0].Changes[0].Value.Messages
	if len(msgs) != 1 || msgs[0].Type != "text" {
		t.Fatalf("messages=%+v", msgs)
	}
	if msgs[0].Text == nil || msgs[0].Text.Body == "" {
		t.Fatal("text body vazio")
	}
}

func TestParseWebhook_StatusFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "webhook_status.json"))
	if err != nil {
		t.Fatal(err)
	}
	ev, err := ParseWebhook(raw)
	if err != nil {
		t.Fatal(err)
	}
	st := ev.Entry[0].Changes[0].Value.Statuses
	if len(st) != 1 || st[0].Status != "delivered" {
		t.Fatalf("statuses=%+v", st)
	}
}
