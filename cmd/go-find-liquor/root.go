package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/toozej/go-find-liquor/internal/runner"
	"github.com/toozej/go-find-liquor/pkg/config"
	"github.com/toozej/go-find-liquor/pkg/man"
	"github.com/toozej/go-find-liquor/pkg/version"
)

var (
	configFile string
	once       bool
)

var rootCmd = &cobra.Command{
	Use:              "go-find-liquor",
	Short:            "Oregon Liquor Search Notification Service",
	Long:             `Oregon Liquor Search Notification Service using the OLCC Liquor Search website, Go, and the nikoksr/notify library`,
	Args:             cobra.ExactArgs(0),
	PersistentPreRun: rootCmdPreRun,
	Run:              rootCmdRun,
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	// Get configuration
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create runner
	r, err := runner.NewRunner(conf)
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Info("Received termination signal, shutting down...")
		cancel()
	}()

	// Run once or continuously
	if once {
		if err := r.RunOnce(ctx); err != nil {
			log.Errorf("Failed to run single search: %v", err)
		}
	} else {
		log.Infof("Starting continuous search with interval %.0f", conf.Interval.Hours())
		if err := r.Start(ctx); err != nil {
			log.Errorf("Failed to run continuous search: %v", err)
		}
	}
}

func rootCmdPreRun(cmd *cobra.Command, args []string) {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return
	}
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// Set custom config file if specified
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Failed to read config file %s: %v", configFile, err)
		}
		log.Infof("Using config file: %s", configFile)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	_, err := maxprocs.Set()
	if err != nil {
		log.Error("Error setting maxprocs: ", err)
	}

	// create rootCmd-level flags
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug-level logging")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file path")
	rootCmd.Flags().BoolVarP(&once, "once", "o", false, "Run search once and exit")

	// add sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)
}
