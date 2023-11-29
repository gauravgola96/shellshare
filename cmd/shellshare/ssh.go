package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"githug.com/gauravgola96/shellshare/pkg/shellshare"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ShellShare is P2P file sharing server",
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetConfigName("default")
		viper.SetConfigType("yaml")
		viper.SetConfigFile("./config/config.yaml")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err).Msgf("Error loading config file: %v", viper.ConfigFileUsed())
			return
		}
		viper.AutomaticEnv()
		err = shellshare.SSHServer()
		if err != nil {
			return
		}
	},
}

func init() {
	shellshareCmd.AddCommand(sshCmd)
}
