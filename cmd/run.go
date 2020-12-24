package run

import (
	"fmt"
	"os"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "capi 😉",
	Run: func(cmd *cobra.Command, args []string) {
		idGenerator := NewIDGenerator()
		idForRootProcess := idGenerator.NewID()

		// validate fields
		logging.WithFields(logrus.Fields{
			"version":          version,
			"build time":       buildTime,
			"event id":         idForRootProcess,
			"config file path": viperConfig.GetString(configFilePathName),
			"log format":       viperConfig.GetString(logFormatName),
			"log level":        viperConfig.GetString(logLevelName),
			"log output":       viperConfig.GetString(logOutputName),
			"syslog tag":       viperConfig.GetString(syslogTagName),

			"consul address":                                    viperConfig.GetString(consulAddressName),
			"consul timeout":                                    viperConfig.GetDuration(consulTimeoutName),
			"is put mode":                                       viperConfig.GetBool(modeIsPutName),
			"is force update mode":                              viperConfig.GetBool(isForceUpdateName),
			"consul root path for lbos cluster":                 viperConfig.GetString(consulRootPathName),
			"folder for json files for send to consul":          viperConfig.GetString(dataDirPathName),
			"manifest key name for service. Example: manifest":  viperConfig.GetString(consulManifestName),
			"folder and consul key-folder name for app servers": viperConfig.GetString(consulAppServersFolderName),
		}).Info("")
		// TODO: validate paths
		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("program running")

		consulKVClient, err := newconsulKVClient(viperConfig.GetString(consulAddressName),
			viperConfig.GetDuration(consulTimeoutName))
		if err != nil {
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("make consul api kv client error: %v",
				err)
		}

		if viperConfig.GetBool(modeIsPutName) {
			if err := putRootPathToConsul(consulKVClient,
				viperConfig.GetString(consulRootPathName)); err != nil {
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("put root path %v to consul error: %v",
					viperConfig.GetString(consulRootPathName),
					err)
			}

			if err := putConsulData(consulKVClient,
				viperConfig.GetString(dataDirPathName),
				viperConfig.GetString(consulManifestName),
				viperConfig.GetString(consulAppServersFolderName),
				viperConfig.GetString(consulRootPathName),
				viperConfig.GetBool(isForceUpdateName)); err != nil {
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("put data error: %v",
					err)
			}
			// TODO: more logs
			logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("put new services success")
		} else {
			consulData, err := readConsulData(consulKVClient,
				viperConfig.GetString(consulRootPathName),
				viperConfig.GetString(consulAppServersFolderName),
				viperConfig.GetString(consulManifestName))
			if err != nil {
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Fatalf("read consul data error: %v",
					err)
			}
			if consulData == nil {
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("no services in consul by path %v",
					viperConfig.GetString(consulRootPathName))
			} else {
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("services: %v",
					consulData)
				logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Infof("total services: %v",
					len(consulData))
			}
		}

		logging.WithFields(logrus.Fields{"event id": idForRootProcess}).Info("program stopped")
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newconsulKVClient(consulAddress string, consulTimeout time.Duration) (*consulapi.KV, error) {
	// Get a new client
	clientAPI := consulapi.DefaultConfig()
	clientAPI.Address = consulAddress
	clientAPI.WaitTime = consulTimeout
	client, err := consulapi.NewClient(clientAPI)
	if err != nil {
		return nil, err
	}
	return client.KV(), nil
}
