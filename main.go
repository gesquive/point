package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var version = "v0.1.0"
var dirty = ""

var cfgFile string

var displayVersion string
var showVersion bool
var verbose bool
var debug bool

var server *Server

func main() {
	displayVersion = fmt.Sprintf("reflect %s%s",
		version,
		dirty)
	Execute(displayVersion)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "reflect",
	Short: "A web client info server",
	Long:  `A web API for client browser information`,
	Run:   run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/%s\n",
		RootCmd.HelpTemplate(), displayVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"Path to a specific config file (default \"./config.yml\")")
	RootCmd.PersistentFlags().StringP("log-file", "l", "",
		"Path to log file (default \"/var/log/reflect.log\")")

	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version number and exit")
	RootCmd.PersistentFlags().StringP("address", "a", "0.0.0.0",
		"The IP address to bind the web server too")
	RootCmd.PersistentFlags().IntP("port", "p", 8080,
		"The port to bind the webserver too")

	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Print logs to stdout instead of file")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Include debug statements in log output")
	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("reflect")
	viper.AutomaticEnv()
	viper.BindEnv("log_file")
	viper.BindEnv("address")
	viper.BindEnv("port")

	viper.BindPFlag("log_file", RootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("web.address", RootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("web.port", RootCmd.PersistentFlags().Lookup("port"))

	viper.SetDefault("log_file", "/var/log/reflect.log")
	viper.SetDefault("web.address", "0.0.0.0")
	viper.SetDefault("web.port", 8080)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath(".")                   // add current directory as first search path
	viper.AddConfigPath("$HOME/.config/reflect") // add home directory to search path
	viper.AddConfigPath("/etc/reflect")          // add etc to search path
	viper.AutomaticEnv()                       // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !showVersion {
			if !strings.Contains(err.Error(), "Not Found") {
				log.Infof("Error opening config: ", err)
			}
		}
	}
}

func run(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Println(displayVersion)
		os.Exit(0)
	}

	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	logFilePath := getLogFilePath(viper.GetString("log_file"))
	log.Debugf("config: log_file=%s", logFilePath)
	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file=%v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Infof("config: file=%s", viper.ConfigFileUsed())

	address := viper.GetString("web.address")
	port := viper.GetInt("web.port")

	// finally, run the webserver
	server := NewServer()
	server.Run(fmt.Sprintf("%s:%d", address, port))
}

func getLogFilePath(defaultPath string) (logPath string) {
	fi, err := os.Stat(defaultPath)
	if err == nil && fi.IsDir() {
		logPath = path.Join(defaultPath, "reflect.log")
	} else {
		logPath = defaultPath
	}
	return
}
