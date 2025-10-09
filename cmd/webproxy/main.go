package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yvv4git/web-proxy/internal/webproxy"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "webproxy",
		Short: "Run application, show help",
	}

	var configPath string
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run web proxy server",
		Run: func(cmd *cobra.Command, args []string) {
			webproxy.RunWebProxy(configPath)
		},
	}

	runCmd.PersistentFlags().
		StringVarP(&configPath, "config", "c", "config.toml", "Path to config file")

	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
