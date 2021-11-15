# Dolphin

![test](https://github.com/ghosind/dolphin/workflows/Test/badge.svg)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/de7fbdc27cd3411b9a2d57d34eae44d2)](https://www.codacy.com/gh/ghosind/dolphin/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ghosind/dolphin&amp;utm_campaign=Badge_Grade)
![Version Badge](https://img.shields.io/github/v/release/ghosind/dolphin)
![License Badge](https://img.shields.io/github/license/ghosind/dolphin)
[![Go Reference](https://pkg.go.dev/badge/github.com/ghosind/dolphin.svg)](https://pkg.go.dev/github.com/ghosind/dolphin)

Dolphin is a simple web framework for Golang.

## Installation

1. Install dolphin by Go cli tool:

    ```bash
    go get -u github.com/ghosind/dolphin
    ```

2. Import dolphin in your code:

    ```go
    import "github.com/ghosind/dolphin"
    ```

## Getting Started

1. The following example shows how to implement a simple web service, and it'll reads parameter name from request query and return the message as a JSON object.

    ```go
    package main

    import (
      "fmt"

      "github.com/ghosind/dolphin"
    )

    func handler(ctx *dolphin.Context) {
      name := ctx.Query("name")

      ctx.JSON(ctx.O{
        "message": fmt.Sprintf("Hello, %s", name),
      })
    }

    func main() {
      app := dolphin.Default()

      app.Use(handler)

      app.Run()
    }
    ```

2. Save the above code as `app.go` and run it with the following command:

    ```bash
    go run app.go
    # Server running at :8080.
    ```

3. Visit `http://locahost:8080?name=dolphin` by your browser or other tools to see the result.

## Examples

### Parameters in query string

```go
func handler(ctx *dolphin.Context) {
  firstName := ctx.Query("firstName")
  lastName := ctx.QueryDefault("lastName", "Doe")

  ctx.JSON(ctx.O{
    "message": fmt.Sprintf("Hello, %s ", firstName, lastName),
  })
}
```

### Parsed request body

```go
type Payload struct {
  Name string `json:"name"`
}

func handler(ctx *dolphin.Context) {
  var payload Payload

  if err := ctx.PostJSON(&payload); err != nil {
    // Return 400 Bad Request
    ctx.Fail(ctx.O{
      "message": "Invalid payload",
    })
    return
  }

  ctx.JSON(ctx.O{
    "message": fmt.Sprintf("Hello, %s ", payload.Name),
  })
}
```

### Custom middleware

```go
func CustomMiddleware() {
  // Do something

  return func (ctx *dolphin.Context) {
    // Do something before following handlers

    ctx.Next()

    // Do something after following handlers
  }
}

func main() {
  app := dolphin.Default()

  app.Use(CustomMiddleware())

  // ...
}
```

## License

Distributed under the MIT License. See LICENSE file for more information.
