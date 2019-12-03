package user

import (
	"fmt"
	"os"

	"github.com/adrienaury/owl/pkg/domain/password"
	"github.com/spf13/cobra"
)

var (
	// local flags
	flagHashAlgorithm    string
	flagCharDomain       string
	flagUpperCaseLetters bool
	flagLowerCaseLetters bool
	flagNumbers          bool
	flagSpecials         bool
	flagLength           uint
)

// initPasswordCommand initialize the cli user password command
func initPasswordCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "new-password [ID]",
		Short:   "Assign a new password to user and send it by mail",
		Long:    "",
		Aliases: []string{"passwd"},
		Example: fmt.Sprintf(`  %[1]s user new-password joker`, parentCmd.Root().Name()),
		Args:    cobra.ExactArgs(1),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]

			domain := password.NewDomain("")
			if flagUpperCaseLetters {
				domain = domain.MergeWith(password.UpperCaseLetters)
			}
			if flagLowerCaseLetters {
				domain = domain.MergeWith(password.LowerCaseLetters)
			}
			if flagNumbers {
				domain = domain.MergeWith(password.Numbers)
			}
			if flagSpecials {
				domain = domain.MergeWith(password.ASCIISpecials)
			}
			if domain.Len() == 0 {
				domain = password.NewDomain(flagCharDomain)
			}

			password, err := passwordDriver.GetRandomPassword(domain, flagLength)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			hash, err := passwordDriver.GetHash(flagHashAlgorithm, password)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if err := userDriver.AssignPassword(id, hash); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			fmt.Println(password, hash)

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Assigned new random password to user '%v' in unit '%v' of realm '%v'.", id, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	cmd.Flags().StringVar(&flagHashAlgorithm, "alg", "SSHA256", "password hash algorithm, one of MD5, SMD5, SHA, SHA1, SSHA, SSHA1, SHA224, SSHA224, SHA256, SSHA256, SHA384, SSHA384, SHA512, SSHA512")
	cmd.Flags().StringVar(&flagCharDomain, "char", string(password.Standard.AsSlice()), "list of characters for the password generation, repeated character appears more often")
	cmd.Flags().BoolVarP(&flagLowerCaseLetters, "lowercase-letters", "l", false, "use lowercase letters for password generation, if set --char is ignored")
	cmd.Flags().BoolVarP(&flagUpperCaseLetters, "uppercase-letters", "L", false, "use uppercase letters for password generation, if set --char is ignored")
	cmd.Flags().BoolVarP(&flagNumbers, "numbers", "n", false, "use numbers for password generation, if set --char is ignored")
	cmd.Flags().BoolVarP(&flagSpecials, "specials", "s", false, "use special characters for password generation, if set --char is ignored")
	cmd.Flags().UintVar(&flagLength, "len", 12, "password length")
	parentCmd.AddCommand(cmd)
}
