package main

import (
	"fmt"
	"os"

	"github.com/whale-team/whaleEcho/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "help"}

func main() {
	rootCmd.AddCommand(cmd.WebSocketCmd, cmd.ClientCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("main: cobra command execution failed, err:%+v\n", err)
		os.Exit(1)
	}
}
