# wa-meta

Pacote Go para protocolo **Meta WhatsApp Cloud API** (library):

- validar assinatura `X-Hub-Signature-256`
- validar challenge de subscribe (retorna o eco para o body HTTP)
- parse de webhook JSON (**sem** servidor HTTP)
- enviar mensagens e baixar mídia via Graph API
- HTTP client injetável (`HTTPDoer`)

O **app consumidor** expõe o endpoint HTTP e responde **200** para evitar reenvio da Meta. Este pacote não abre porta e não envia “ACK” sozinho.

```go
import "github.com/rafapasa/wa-meta"

// GET subscribe
echo, err := wameta.ValidateChallenge(verifyToken, mode, token, challenge)
// POST webhook
if err := wameta.ValidateSignature(appSecret, sigHeader, body); err != nil { ... }
ev, err := wameta.ParseWebhook(body)

client := wameta.NewClient(wameta.Config{...})
client.SendText(ctx, wameta.SendTextInput{To: "...", Body: "oi"})
```

```bash
go test ./...
```
