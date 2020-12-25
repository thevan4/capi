package run

import (
	"fmt"
	"os"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	logger "github.com/thevan4/logrus-wrapper"
)

var rootCmd = &cobra.Command{
	Use:   "capi",
	Short: "consul api worker for lbos 😉",
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var getConsulRequest = &cobra.Command{
	Use:   "get",
	Short: "get services at consul root path",
	RunE: func(cmd *cobra.Command, args []string) error {
		consulKVClient, logging, idForRootProcess := initForStart()
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("start get services")
		consulData, err := readConsulData(consulKVClient,
			viper.GetString("consul-root-path"),
			viper.GetString("app-servers-folder"),
			viper.GetString("manifest-name"))
		if err != nil {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("read consul data error: %v",
				err)
		}
		if consulData == nil {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("no services in consul by path %v",
				viper.GetString("consul-root-path"))
		} else {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("services: %v",
				consulData)
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("total services: %v",
				len(consulData))
		}
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("get services success")
		return nil
	},
}

var putConsulRequest = &cobra.Command{
	Use:   "put",
	Short: "put service in consul",
	RunE: func(cmd *cobra.Command, args []string) error {
		consulKVClient, logging, idForRootProcess := initForStart()
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("start put services")
		if err := putRootPathToConsul(consulKVClient,
			viper.GetString("consul-root-path")); err != nil {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("put root path %v to consul error: %v",
				viper.GetString("consul-root-path"),
				err)
		}
		if err := putConsulData(consulKVClient,
			viper.GetString("data-dir-path-names"),
			viper.GetString("manifest-name"),
			viper.GetString("app-servers-folder"),
			viper.GetString("consul-root-path"),
			viper.GetBool("force-update-keys")); err != nil {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("put data error: %v",
				err)
		}
		// TODO: more logs
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("put new services success")
		return nil
	},
}

func initForStart() (*consulapi.KV, *logrus.Logger, string) {
	// init logs
	newLogger := &logger.Logger{
		Output:           []string{viper.GetString("log-output")},
		Level:            viper.GetString("log-level"),
		Formatter:        viper.GetString("log-format"),
		LogEventLocation: viper.GetBool("log-event-location"),
	}
	logging, err := logger.NewLogrusLogger(newLogger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	idGenerator := NewIDGenerator()
	idForRootProcess := idGenerator.NewID()

	// validate fields
	logging.WithFields(logrus.Fields{
		"version":          version,
		"build time":       buildTime,
		"event id":         idForRootProcess,
		"config file path": viper.GetString("config"),
		"log format":       viper.GetString("log-format"),
		"log level":        viper.GetString("log-level"),
		"log output":       viper.GetString("log-output"),
		"syslog tag":       viper.GetString("syslog-tag"),

		"consul address":                                    viper.GetString("consul-address"),
		"consul timeout":                                    viper.GetDuration("consul-timeout"),
		"is put mode":                                       viper.GetBool("mode"),
		"is force update mode":                              viper.GetBool("force-update-keys"),
		"consul root path for lbos cluster":                 viper.GetString("consul-root-path"),
		"folder for json files for send to consul":          viper.GetString("data-dir-path-names"),
		"manifest key name for service. Example: manifest":  viper.GetString("manifest-name"),
		"folder and consul key-folder name for app servers": viper.GetString("app-servers-folder"),
	}).Info("")
	// TODO: validate paths

	consulKVClient, err := newconsulKVClient(viper.GetString("consul-address"),
		viper.GetDuration("consul-timeout"))
	if err != nil {
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("make consul api kv client error: %v",
			err)
	}
	return consulKVClient, logging, idForRootProcess
}
