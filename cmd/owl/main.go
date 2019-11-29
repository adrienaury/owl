package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/adrienaury/owl/cmd/owl/realm"
	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/cmd/owl/unit"
	"github.com/adrienaury/owl/cmd/owl/user"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version   string
	commit    string
	buildDate string
	builtBy   string

	// session
	home          = initHome()
	globalSession = session.NewSession(path.Join(home, "session.yaml"))

	// global flags
	flagCfgFile string
	flagRealm   string
	flagUnit    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "owl [action]",
	Short: "Command line tools for managing realms of units, users and groups",
	Long:  `Owl is a tools to manage realms of units, users and groups.`,
	Example: `  owl realm set dev ldap://dev.my-company.com/dc=example,dc=com
  owl realm login dev
  owl unit apply <<< '{"id": "my-unit"}'
  owl unit use my-unit
  owl user create <<< '{"id": "batman", "first-name": ["Bruce"], "last-name": ["Wayne"]}'
  owl user apply joker first-name=Arthur last-name=Flake email=arthur.flake@gotham.dc
  owl user list -o table`,
	Version: fmt.Sprintf("%v (commit=%v date=%v by=%v)", version, commit, buildDate, builtBy),
}

func main() {
	defer globalSession.Dump()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags
	rootCmd.PersistentFlags().StringVar(&flagCfgFile, "config", "", "config file (default is $HOME/.owl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "target realm")
	rootCmd.PersistentFlags().StringVar(&flagUnit, "unit", "", "target unit")

	rootCmd.AddCommand(realm.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
	rootCmd.AddCommand(unit.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
	rootCmd.AddCommand(user.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
}

func initConfig() {
	// restore session
	if _, err := globalSession.Restore(); err != nil {
		// TODO logger
		fmt.Println(err.Error())
	}

	if flagCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagCfgFile)
	} else {
		// Search config in home directory.
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// TODO: logger
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if strings.TrimSpace(flagUnit) == "" {
		flagUnit = globalSession.Unit
	}

	if strings.TrimSpace(flagRealm) == "" {
		flagRealm = globalSession.Realm
	}

	backend := newBackend()
	backend.SetUnit(flagUnit)
	credentialsDriver := newCredentialsDriver(backend)
	realmDriver := newRealmDriver()

	r, _ := realmDriver.Get(flagRealm)
	if r != nil {
		c, _ := credentialsDriver.Get(r.URL(), r.Username())
		if c != nil {
			backend.SetCredentials(c)
		}
	}

	realm.SetDrivers(realmDriver, credentialsDriver)
	realm.SetSession(globalSession)

	unitDriver := newUnitDriver(backend)
	unit.SetDrivers(unitDriver, realmDriver, credentialsDriver)
	unit.SetSession(globalSession)

	userDriver := newUserDriver(backend)
	user.SetDrivers(userDriver, realmDriver, credentialsDriver)
}

func initHome() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Mkdir(home+"/.owl", 0644)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	return path.Join(home, ".owl")
}
