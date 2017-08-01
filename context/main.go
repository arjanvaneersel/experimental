package main

import (
	"context"
	"time"
	"fmt"
	"bufio"
	"os"
	"log"
)

func withTimeOut() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	sleepAndTalk(ctx, 5 * time.Second, "Hello gopher!")
}

func WithCancel() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		cancel()
	}()

	sleepAndTalk(ctx, 5 * time.Second, "Hello gopher!")
}

func Background() {
	ctx := context.Background()
	sleepAndTalk(ctx, 5 * time.Second, "Hello gopher!")
}

func sleepAndTalk(ctx context.Context, d time.Duration, s string) {
	select {
	case <-time.After(d):
		fmt.Println(s)
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}

func main() {
	Background()
	WithCancel()
	withTimeOut()
}
