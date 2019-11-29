package main

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/cmd/owl/realm"
	"github.com/adrienaury/owl/cmd/owl/unit"
	"github.com/adrienaury/owl/cmd/owl/user"
	"github.com/spf13/cobra"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version   string
	commit    string
	buildDate string
	builtBy   string

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "target realm")
	rootCmd.PersistentFlags().StringVar(&flagUnit, "unit", "", "target unit")

	rootCmd.AddCommand(realm.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
	rootCmd.AddCommand(unit.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
	rootCmd.AddCommand(user.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
}

func initConfig() {
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

	unitDriver := newUnitDriver(backend)
	unit.SetDrivers(unitDriver, realmDriver, credentialsDriver)

	userDriver := newUserDriver(backend)
	user.SetDrivers(userDriver, realmDriver, credentialsDriver)
}
