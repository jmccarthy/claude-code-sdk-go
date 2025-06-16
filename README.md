# Claude Code SDK for Go

This repository contains an experimental Go implementation of the Claude Code SDK.
The goal is feature parity with the existing Python SDK while providing an idiomatic
Go interface.

**Status:** Work in progress. See `PORT_TODO.md` for the current porting plan and
completed tasks.

## Installation

Go 1.20 or newer is required.

```bash
go get github.com/anthropics/claude-code-sdk-go/...
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/anthropics/claude-code-sdk-go/claudecode"
)

func main() {
    ctx := context.Background()
    ch, errCh, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
    if err != nil {
        log.Fatal(err)
    }
    for {
        select {
        case msg, ok := <-ch:
            if !ok {
                ch = nil
                continue
            }
            log.Printf("%+v", msg)
        case e, ok := <-errCh:
            if ok {
                log.Printf("error: %v", e)
            }
            errCh = nil
        }
        if ch == nil && errCh == nil {
            break
        }
    }
}
```

## Examples

See `examples/quick_start.go` for a more complete example of using the package.

## License

MIT
