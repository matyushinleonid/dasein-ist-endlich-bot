package main

import (
	"context"
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/listener"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/notifier"
	"github.com/spf13/cobra"
)

var (
	rootCmd     = &cobra.Command{Use: "dasein-ist-endlich-bot"}
	listenerCmd = &cobra.Command{
		Use:  "listener",
		RunE: runBot("listener", func(b *core.DaseinBot) BotRunner { return listener.New(b) }),
	}
	notifierCmd = &cobra.Command{
		Use:  "notifier",
		RunE: runBot("notifier", func(b *core.DaseinBot) BotRunner { return notifier.New(b) }),
	}
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "./configs/local.yaml", "Path to the config file")
	rootCmd.AddCommand(listenerCmd, notifierCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type BotRunner interface {
	Run(ctx context.Context) error
}

func runBot(name string, factory func(*core.DaseinBot) BotRunner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		rootLogger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)).WithName(name)
		ctx = logr.NewContext(ctx, rootLogger)
		logger := logr.FromContextOrDiscard(ctx)

		conf, err := config.NewConfig(configPath)
		if err != nil {
			logger.Error(err, "loading config", "path", configPath)
			return err
		}
		logger.Info("configuration loaded", "path", configPath)

		bot := core.New(conf)
		return factory(bot).Run(ctx)
	}
}
