# Go Lumina

This project provides a SDK for interacting with a **Large Language Model (LLM)**.

## Required

Version:
- Golang (Go): 1.25.4

## Instruction manual

1) Import this library in your project:
```golang
go get -u "https://github.com/nexula-rg/go-lumina"
```

2) Create new config:
```golang
import (
    "github.com/nexula-rg/go-lumina/lumina"
)

client, err := lumina.NewClient()

// Data description
var client *lumina.Client   // Link on Client
var err error            // Client creating error
```

3) Make request:
```golang
import (
    "github.com/nexula-rg/go-lumina/lumina"
)

resp, err := client.MakeRequest("Question")

// Data description
var resp lumina.Response   // JSON answer from API LLM

// Response structure description
type Response struct {
     Answer string `json:"answer"`
}

var err error             // Request error
```

4) Make request with context from package "context":
```golang
import (
    "context"

    "github.com/nexula-rg/go-lumina/lumina"
)

ctx, cancel := context(context.Background(), time.Second*5)
defer cancel()

resp, err := client.MakeRequest(ctx, "Question")

// Data description
var resp lumina.Response   // JSON answer from API LLM

// Response structure description
type Response struct {
	Answer string `json:"answer"`
}

var err error             // Request error
```

5) Make async request
```golang
import (
    "github.com/nexula-rg/go-lumina/lumina"
)

resp, err := client.MakeRequest("Question")

// Data description
var resp lumina.AsyncResponse   // JSON answer from API LLM

// Async Response structure description
type AsyncResponse struct {
    AnswerUUID string `json:"answer_uuid"`
    Status     string `json:"status"`
}

var err error             // Request error
```

6) Error handling
```golang
import (
    "log"

    "github.com/nexula-rg/go-lumina/lumina"
)

answer, err := client.MakeRequest("Question")
if err != nil {
    if errResp, ok := err.(*lumina.ErrorResponse); ok {
        log.Printf("error: %s, status code: %d", errResp.Error(), errResp.StatusCode)
    } else {
        log.Println("unknown error:", err)
    }
}

// Error structure description
type ErrorResponse struct {
    StatusCode int    `json:"status_code"`
    Msg        string `json:"error"`
}
```

7) Set env variable at terminal in your project:
```bash
export GO_LUMINA_URL="http://localhost:8080"
```

or on Windows

```cmd
set GO_LUMINA_URL=http://localhost:8080
```
