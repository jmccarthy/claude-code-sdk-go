package main

import (
	"context"
	"log"

	"github.com/anthropics/claude-code-sdk-go/claudecode"
)

func main() {
	ctx := context.Background()

	ch, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range ch {
		log.Printf("%+v", msg)
	}
}
