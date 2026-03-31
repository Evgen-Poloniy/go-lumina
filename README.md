# Go Lumina

**Go Lumina** is a Go SDK for interacting with a **Large Language Model (LLM)** via HTTP API.

## Requirements

- Golang (Go) version **1.25.4** or higher
- LLM service URL set via environment variable `LUMINA_URL`

## Installation

Import the library into your project:

```bash
go get -u github.com/nexula-rg/go-lumina
```

## Usage

### 1. Create a Client

```golang
import (
    "github.com/nexula-rg/go-lumina/lumina"
)

client, err := lumina.NewClient("Required API-KEY")
if err != nil {
    panic(err)
}

// client is the Client instance to interact with the LLM API
```

### 2. Make a Request to the LLM

```golang
answer, err := client.MakeRequest("Question")
if err != nil {
    panic(err)
}

fmt.Println("LLM answer:", answer)
```

**Response scheme:**

```golang
var answer string
```

### 3. Request with Context (`context.Context`)

```golang
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.MakeRequestCtx(ctx, "Question")
if err != nil {
    panic(err)
}

fmt.Println("LLM answer:", resp)
```

> Or use any context whatever is required
> Examples: timeout, cancel, application graceful shutdown cancel context

### 4. Check API Availability (Ping)

```golang
if err := client.Ping(); err != nil {
    fmt.Println("LLM service is unavailable:", err)
} else {
    fmt.Println("LLM service is available")
}
```

With context:

```golang
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

if err := client.PingCtx(ctx); err != nil {
    fmt.Println("LLM service is unavailable:", err)
} else {
    fmt.Println("LLM service is available")
}
```

### 5. Error Handling

```golang
answer, err := client.MakeRequest("Question")
if err != nil {
    if errResp, ok := err.(*lumina.ErrorResponse); ok {
        fmt.Printf("Error: %s: %s, status code: %d\n", errResp.Code, errResp.Error(), errResp.StatusCode)
    } else {
        fmt.Println("Unknown error:", err)
    }
}
```

**Error structure:**

```golang
type ErrorResponse struct {
	StatusCode int                      // HTTP status code
	Details    struct {
		Code    string `json:"code"`    // Code for frontend
		Message string `json:"message"` // Error message
	} `json:"error"`
}
```

### 6. Set API URL

On Linux/macOS:

```bash
export LUMINA_URL="http://localhost:8080"
```

On Windows:

```cmd
set LUMINA_URL=http://localhost:8080
```
