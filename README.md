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
    for msg := range ch {
        log.Printf("%+v", msg)
    }
    if e := <-errCh; e != nil {
        log.Fatalf("query error: %v", e)
    }
}
```

## Examples

See `examples/quick_start.go` for a more complete example of using the package.

### Custom CLI Path

If the `claude` CLI binary is not on your `PATH`, set `Options.CLIPath` to the
location of the executable when calling `Query`.

## License

MIT
