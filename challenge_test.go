package wameta

import "testing"

func TestValidateChallenge_OK(t *testing.T) {
	echo, err := ValidateChallenge("tok", "subscribe", "tok", "12345")
	if err != nil {
		t.Fatal(err)
	}
	if echo != "12345" {
		t.Fatalf("echo=%s", echo)
	}
}

func TestValidateChallenge_BadToken(t *testing.T) {
	_, err := ValidateChallenge("tok", "subscribe", "other", "12345")
	if err == nil {
		t.Fatal("esperava erro")
	}
}
