package main

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/cmd/owl/realm"
	"github.com/spf13/cobra"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version   string
	commit    string
	buildDate string
	builtBy   string
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

	rootCmd.AddCommand(realm.NewCommand("owl", os.Stderr, os.Stdout, os.Stdin))
}

func initConfig() {
	realm.SetDrivers(realmDriver(), credentialsDriver())
}
