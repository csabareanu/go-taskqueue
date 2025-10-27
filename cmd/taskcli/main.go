package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/csabareanu/taskqueue/pkg/config"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded Configuration:\n")
	fmt.Println("---------------------")
	fmt.Printf("Server Address : %s\n", cfg.ServerAddress)
	fmt.Printf("Max Workers    : %d\n", cfg.MaxWorkers)
	fmt.Printf("Log Level      : %s\n", cfg.LogLevel)
	fmt.Println("---------------------")
}
