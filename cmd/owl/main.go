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
	flagRealm string
	flagUnit  string
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
	if _, err := globalSession.Restore(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer globalSession.Dump()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// fix standards streams
	rootCmd.SetIn(os.Stdin)
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	// global flags
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "target realm")
	rootCmd.PersistentFlags().StringVar(&flagUnit, "unit", "", "target unit")

	realm.InitCommand(rootCmd)
	unit.InitCommand(rootCmd)
	user.InitCommand(rootCmd)
}

func initConfig() {
	if strings.TrimSpace(flagUnit) == "" {
		flagUnit = globalSession.Unit
	}

	if strings.TrimSpace(flagRealm) == "" {
		flagRealm = globalSession.Realm
	}

	backend := newBackend()
	credentialsDriver := newCredentialsDriver(&backend)
	realmDriver := newRealmDriver()
	unitDriver := newUnitDriver(&backend)
	userDriver := newUserDriver(&backend)

	realm.SetDrivers(realmDriver, credentialsDriver)
	realm.SetSession(globalSession)

	unit.SetDrivers(unitDriver, realmDriver, credentialsDriver)
	unit.SetSession(globalSession)

	user.SetDrivers(userDriver, unitDriver, realmDriver, credentialsDriver)
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
