package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "chirpstack-simulator",
	Short: "ChirpStack Simulator CLI",
	Long: `ChirpStack Simulator is a CLI tool for simulating ChirpStack devices.
	> documentation & support: https://www.chirpstack.io/`,
	Version: version,
	RunE:    run,
}
