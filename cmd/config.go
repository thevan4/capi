package run

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// For builds with ldflags
var (
	version   = "unknown"
	buildTime = "unknown"
	// 	commit  = "TBD @ ldflags"
	// 	branch  = "TBD @ ldflags"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config",
		"c",
		"./capi.yaml",
		"Path to config file. Example value: './capi.yaml'")
	rootCmd.PersistentFlags().String("log-output",
		"stdout",
		"Log output. Example values: 'stdout', 'syslog'")
	rootCmd.PersistentFlags().String("log-level",
		"trace", "Log level. Example values: 'info', 'debug', 'trace'")
	rootCmd.PersistentFlags().String("log-format",
		"text",
		"Log format. Example values: 'text', 'json'")
	rootCmd.PersistentFlags().String("syslog-tag",
		"",
		"Syslog tag. Example: 'capi-worker'")
	rootCmd.PersistentFlags().Bool("log-event-location",
		true,
		"Log event location (like python)")

	rootCmd.PersistentFlags().StringP("consul-address",
		"a",
		"127.0.0.1:8500",
		"consul address")
	rootCmd.Flags().DurationP("consul-timeout",
		"t",
		2*time.Second,
		"consul timeout")
	rootCmd.PersistentFlags().StringP("key-for-remove",
		"k",
		"lbos/t1-cluster-0/",
		"key for remove. Also can remove key tree, when set force key")
	rootCmd.PersistentFlags().BoolP("force-keys",
		"f",
		false,
		"force update/remove keys (bool)")

	rootCmd.PersistentFlags().StringP("consul-root-path",
		"r",
		"lbos/t1-cluster-0/",
		"consul root path for lbos cluster. Example: lbos/t1-cluster-0/")
	rootCmd.PersistentFlags().StringP("data-dir-path-names",
		"d",
		"./json-services",
		"folder for json files for send to consul. Example: ./json-services")
	rootCmd.PersistentFlags().StringP("manifest-name",
		"m",
		"manifest",
		"manifest key name for service. Example: manifest")
	rootCmd.PersistentFlags().StringP("app-servers-folder",
		"s",
		"app-servers",
		"folder and consul key-folder name for app servers. Example: app-servers")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(getConsulRequest, putConsulRequest, delConsulRequest)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("can't read config from file, error:", err)
	}
}
