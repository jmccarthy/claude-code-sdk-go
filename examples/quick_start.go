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
