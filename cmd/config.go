package run

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	logger "github.com/thevan4/logrus-wrapper"
)

// Default values
const (
	defaultConfigFilePath   = "./capi.yaml"
	defaultLogOutput        = "stdout"
	defaultLogLevel         = "trace"
	defaultLogFormat        = "text"
	defaultSystemLogTag     = ""
	defaultLogEventLocation = true

	defaultConsulAddress     = "127.0.0.1:8500"
	defaultConsulTimeout     = 2 * time.Second
	defaultIsPutMode         = false
	defaultIsForceUpdateName = false

	defaultConsulRootPath         = "lbos/t1-cluster-0/"
	defaultDataDirPath            = "./json-services"
	defaultConsulManifest         = "manifest"
	defaultConsulAppServersFolder = "app-servers"
)

// Config names
const (
	configFilePathName   = "config-file-path"
	logOutputName        = "log-output"
	logLevelName         = "log-level"
	logFormatName        = "log-format"
	syslogTagName        = "syslog-tag"
	logEventLocationName = "log-event-location"

	consulAddressName = "consul-address"
	consulTimeoutName = "consul-timeout"
	modeIsPutName     = "mode"
	isForceUpdateName = "force-update-keys"

	consulRootPathName         = "consul-root-path"
	dataDirPathName            = "data-dir-path-names"
	consulManifestName         = "manifest-name"
	consulAppServersFolderName = "app-servers-folder"
)

// // For builds with ldflags
var (
	version   = "unknown"
	buildTime = "unknown"
	// 	commit  = "TBD @ ldflags"
	// 	branch  = "TBD @ ldflags"
)

var (
	viperConfig *viper.Viper
	logging     *logrus.Logger
)

func init() {
	var err error
	viperConfig = viper.New()
	// work with env
	viperConfig.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viperConfig.AutomaticEnv()

	// work with flags
	pflag.StringP(configFilePathName, "c", defaultConfigFilePath, "Path to config file. Example value: './capi.yaml'")
	pflag.String(logOutputName, defaultLogOutput, "Log output. Example values: 'stdout', 'syslog'")
	pflag.String(logLevelName, defaultLogLevel, "Log level. Example values: 'info', 'debug', 'trace'")
	pflag.String(logFormatName, defaultLogFormat, "Log format. Example values: 'text', 'json'")
	pflag.String(syslogTagName, defaultSystemLogTag, "Syslog tag. Example: 'trac-dgen'")
	pflag.Bool(logEventLocationName, defaultLogEventLocation, "Log event location (like python)")

	pflag.StringP(consulAddressName, "a", defaultConsulAddress, "consul address")
	pflag.DurationP(consulTimeoutName, "t", defaultConsulTimeout, "consul timeout")
	pflag.BoolP(modeIsPutName, "p", defaultIsPutMode, "mode: put (true) or read (false)")
	pflag.BoolP(isForceUpdateName, "f", defaultIsForceUpdateName, "force update keys (bool)")

	pflag.StringP(consulRootPathName, "r", defaultConsulRootPath, "consul root path for lbos cluster. Example: lbos/t1-cluster-0/")
	pflag.StringP(dataDirPathName, "d", defaultDataDirPath, "folder for json files for send to consul. Example: ./json-services")
	pflag.StringP(consulManifestName, "m", defaultConsulManifest, "manifest key name for service. Example: manifest")
	pflag.StringP(consulAppServersFolderName, "s", defaultConsulAppServersFolder, "folder and consul key-folder name for app servers. Example: app-servers")

	pflag.Parse()
	if err := viperConfig.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// work with config file
	viperConfig.SetConfigFile(viperConfig.GetString(configFilePathName))
	if err := viperConfig.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	// init logs
	newLogger := &logger.Logger{
		Output:           []string{viperConfig.GetString(logOutputName)},
		Level:            viperConfig.GetString(logLevelName),
		Formatter:        viperConfig.GetString(logFormatName),
		LogEventLocation: viperConfig.GetBool(logEventLocationName),
	}
	logging, err = logger.NewLogrusLogger(newLogger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
