package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	buildVersion = "v0.1.4-dev"
	buildCommit  = ""
	buildDate    = ""
)

var cfgFile string

var showVersion bool
var debug bool

var server *Server

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:    "reflect",
	Short:  "A web client info server",
	Long:   `A web API for client browser information`,
	PreRun: preRun,
	Run:    run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/reflect %s\n",
		RootCmd.HelpTemplate(), buildVersion))
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
		"Display the version info and exit")
	RootCmd.PersistentFlags().StringP("web-address", "a", "0.0.0.0",
		"The IP address to bind the web server too")
	RootCmd.PersistentFlags().IntP("web-port", "p", 2626,
		"The port to bind the webserver too")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Include debug statements in log output")
	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("reflect")
	viper.AutomaticEnv()
	viper.BindEnv("config")
	viper.BindEnv("log_file")
	viper.BindEnv("web_address")
	viper.BindEnv("web_port")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("log_file", RootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("web.address", RootCmd.PersistentFlags().Lookup("web-address"))
	viper.BindPFlag("web.port", RootCmd.PersistentFlags().Lookup("web-port"))

	viper.SetDefault("log_file", "/var/log/reflect.log")
	viper.SetDefault("web.address", "0.0.0.0")
	viper.SetDefault("web.port", 2626)

	dotReplacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(dotReplacer)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")                // name of config file (without extension)
		viper.AddConfigPath(".")                     // add current directory as first search path
		viper.AddConfigPath("$HOME/.config/reflect") // add home directory to search path
		viper.AddConfigPath("/etc/reflect")          // add etc to search path
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !showVersion {
			if !strings.Contains(err.Error(), "Not Found") {
				fmt.Printf("Error opening config: %s\n", err)
			}
		}
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/reflect\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
}

func run(cmd *cobra.Command, args []string) {
	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Infof("running reflect %s", buildVersion)
	if len(buildCommit) > 6 {
		log.Debugf("build: commit=%s", buildCommit[:7])
	}
	if buildDate != "" {
		log.Debugf("build: date=%s", buildDate)
	}
	log.Debugf("build: info=%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)


	logFilePath := getLogFilePath(viper.GetString("log_file"))
	log.Debugf("config: log_file=%s", logFilePath)
	if strings.ToLower(logFilePath) == "stdout" || logFilePath == "-" || logFilePath == "" {
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
