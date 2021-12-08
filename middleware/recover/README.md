# Dolphin Recover Middleware

Recover is a Dolphin framework middleware that provides a recover mechanism for your application. You can set the recover handler to handle the errors, or it'll return a 500 (Internal Server Error) response if the handler is not set.

## Getting Started

```go
import (
  "github.com/ghosind/dolphin"
  "github.com/ghosind/dolphin/middleware/recover"
)

func main() {
  app := dolphin.New()

  app.Use(recover.Recovery())

  app.Run()
}
```

## API

- `Recover(config ...Config) dolphin.HandlerFunc`

  Creates and returns a new recover middleware.

## Config

| Field | Type | Description |
|:------:|:----:|:------------|
| `Handler` | `func (ctx *dolphin.Context, err error)` | Recover handler. |
