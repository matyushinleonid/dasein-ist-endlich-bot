package main

import (
	"context"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/primary"
	"github.com/spf13/cobra"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

var (
	rootCmd = &cobra.Command{
		Use: "dasein-ist-endlich-bot",
	}
	botCmd = &cobra.Command{
		Use: "main-bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			rootLogger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)).WithName("Bot")
			ctx = logr.NewContext(ctx, rootLogger)
			logger := logr.FromContextOrDiscard(ctx)

			conf, err := config.NewConfig(configPath)
			if err != nil {
				logger.Error(err, "loading config", "path", configPath)
				return err
			}
			logger.Info("configuration loaded", "path", configPath)

			b := primary.New(conf)
			return b.Run(ctx)
		},
	}
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "./configs/local.yaml", "Path to the config file")
	rootCmd.AddCommand(botCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
