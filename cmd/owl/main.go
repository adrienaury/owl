package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/adrienaury/owl/cmd/owl/export"
	"github.com/adrienaury/owl/cmd/owl/group"
	"github.com/adrienaury/owl/cmd/owl/realm"
	"github.com/adrienaury/owl/cmd/owl/session"
	"github.com/adrienaury/owl/cmd/owl/unit"
	"github.com/adrienaury/owl/cmd/owl/user"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
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
	flagOutput string
	flagRealm  string
	flagUnit   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "owl [action]",
	Short: "Command line tools for managing realms of units, users and groups",
	Long:  `Owl is a tools to manage realms of units, users and groups.`,
	Example: `  owl realm set dev ldap://dev.my-company.com/dc=my-company,dc=com cn=admin,dc=my-company,dc=com
  owl realm login dev
  owl unit create <<< '{"ID": "my-unit"}'
  owl unit use my-unit
  owl user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'
  owl user ls -o=json | owl import --realm=prod`,
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

	defaultOutput := "json"
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		defaultOutput = "table"
	}

	// global flags
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", defaultOutput, "output format, one of json|yaml|table")
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "target realm")
	rootCmd.PersistentFlags().StringVar(&flagUnit, "unit", "", "target unit")

	realm.InitCommand(rootCmd)
	unit.InitCommand(rootCmd)
	user.InitCommand(rootCmd)
	group.InitCommand(rootCmd)
	export.InitCommand(rootCmd)
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
	passwordDriver := newPasswordDriver(&backend)
	groupDriver := newGroupDriver(&backend)

	realm.SetDrivers(realmDriver, credentialsDriver)
	unit.SetDrivers(unitDriver, realmDriver, credentialsDriver)
	user.SetDrivers(userDriver, passwordDriver, unitDriver, realmDriver, credentialsDriver)
	group.SetDrivers(groupDriver, userDriver, unitDriver, realmDriver, credentialsDriver)
	export.SetDrivers(groupDriver, userDriver, unitDriver, realmDriver, credentialsDriver)

	realm.SetSession(globalSession)
	unit.SetSession(globalSession)
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
