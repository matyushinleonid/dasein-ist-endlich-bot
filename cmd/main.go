package main

import (
	"fmt"
	"os"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "root",
}

var botCmd = &cobra.Command{
	Use: "bot",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfig(configPath)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		bot := bot.NewBot(conf)
		bot.Run()
	},
}

var configPath string

func main() {
	fmt.Println("Starting the application...")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "/etc/config.yaml", "Path to the config file")
	rootCmd.AddCommand(botCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
