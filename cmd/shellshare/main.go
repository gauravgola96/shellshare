package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var shellshareCmd = &cobra.Command{
	Use:   "shellshare",
	Short: "ShellShare is P2P file sharing server",
}

func main() {
	err := shellshareCmd.Execute()
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
		return
	}
}
