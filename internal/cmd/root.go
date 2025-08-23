package cmd

import (
	"fmt"
	"github.com/bart-lute/addigy-tools/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
)

var (
	cfgFile           string
	defaultConfigPath = fmt.Sprintf("%s/.config/addigy-tools", os.Getenv("HOME"))

	rootCmd = &cobra.Command{
		Use:     "addigy-tools",
		Short:   "Useful commands, using the Addigy API",
		Version: config.Version,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s/config.yaml)", defaultConfigPath))

	// ADE Subcommands
	adeListCmd.Flags().BoolP("broken-only", "", false, "Only list broken Automatic Device Enrollments")
	adeCmd.AddCommand(adeListCmd)

	// Policies Subcommands
	werkplekProCmd.AddCommand(werkplekProClientsCmd)

	rootCmd.AddCommand(adeCmd)
	rootCmd.AddCommand(werkplekProCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(defaultConfigPath)
		viper.AddConfigPath(wd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}
	slog.Debug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))

}
