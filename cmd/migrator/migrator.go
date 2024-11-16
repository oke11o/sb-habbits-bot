package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/oke11o/sb-habits-bot/internal/bootstrap"
)

var (
	Version = "dev"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	appname := "sbhabits"
	err := bootstrap.RunMigrator(ctx, os.Args, appname, Version)
	if err != nil {
		fmt.Printf("\nSTOP with error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}
