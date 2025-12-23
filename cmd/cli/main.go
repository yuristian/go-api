package main

import (
	"fmt"
	"os"

	"github.com/yuristian/go-api/internal/cli"
)

func main() {
	args := os.Args

	// Expected pattern:
	// go run ./cmd/cli module create <module_name>
	if len(args) < 4 {
		printUsage()
		os.Exit(1)
	}

	resource := args[1]
	action := args[2]
	name := args[3]

	switch resource {
	case "module":
		handleModuleCommand(action, name)
	default:
		fmt.Println("Unknown resource:", resource)
		printUsage()
		os.Exit(1)
	}
}

func handleModuleCommand(action, name string) {
	switch action {
	case "create":
		if err := cli.GenerateModule(name); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "remove":
		fmt.Println("[CLI] Removing module:", name)
		cli.RemoveModule(name)

	default:
		fmt.Println("Unknown module action:", action)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:

  go run ./cmd/cli module create <module_name>

Examples:

  go run ./cmd/cli module create user
  go run ./cmd/cli module create product
`)
}
