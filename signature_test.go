package wameta

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestValidateSignature_OK(t *testing.T) {
	secret := "testsecret"
	body := []byte(`{"object":"whatsapp_business_account"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if err := ValidateSignature(secret, sig, body); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSignature_Bad(t *testing.T) {
	err := ValidateSignature("x", "sha256=deadbeef", []byte("{}"))
	if err == nil {
		t.Fatal("esperava erro")
	}
}
