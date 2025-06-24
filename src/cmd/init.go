package cmd

import (
	"bytes"
	"os"

	"github.com/deepshore/chirpstack-simulator/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string
var cfgFile string

// Execute is the main entry point for the CLI application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.toml)")
}

func initConfig() {
	config.Version = version
	// Handle the config file path
	if cfgFile == "" {
		cfgFile = "config.toml"
	}

	log.Infof("Using config file: %s", cfgFile)

	// Read the configuration file
	b, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Configure viper to read the config file
	viper.SetConfigType("toml")
	if err := viper.ReadConfig(bytes.NewReader(b)); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	if err := viper.Unmarshal(&config.C); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	log.Infof("Configuration loaded successfully")
}
