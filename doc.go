// Package wameta implements Meta WhatsApp Cloud API protocol helpers:
// signature verification, webhook challenge, webhook JSON parse, send message and media download.
//
// This package is a library: it does not start an HTTP server. The host application
// owns the HTTP endpoint and must respond 200 to Meta to avoid retries.
package wameta
