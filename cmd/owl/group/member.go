package group

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initMemberCommand initialize the cli group member command
func initMemberCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "member {add,remove}",
		Short:   "Manage group members",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s group member add my-group batman`, parentCmd.Root().Name()),
	}
	parentCmd.AddCommand(cmd)
	initMemberAddCommand(cmd)
	initMemberRemoveCommand(cmd)
}

// initMemberAddCommand initialize the cli member add command
func initMemberAddCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add member to group",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s group member add my-group batman ...`, parentCmd.Root().Name()),
		Args:    cobra.MinimumNArgs(2),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			groupID := args[0]

			err := groupDriver.AddMembers(groupID, args[1:]...)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Modified group '%v' in unit '%v' of realm '%v'.", groupID, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}

// initMemberRemoveCommand initialize the cli member remove command
func initMemberRemoveCommand(parentCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove member from group",
		Long:    "",
		Example: fmt.Sprintf(`  %[1]s group member remove my-group batman ...`, parentCmd.Root().Name()),
		Args:    cobra.MinimumNArgs(2),
		PreRun:  initCredentialsAndUnit,
		Run: func(cmd *cobra.Command, args []string) {
			groupID := args[0]

			err := groupDriver.RemoveMembers(groupID, args[1:]...)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			flagUnit := cmd.Flag("unit")
			flagRealm := cmd.Flag("realm")

			cmd.PrintErrf("Modified group '%v' in unit '%v' of realm '%v'.", groupID, flagUnit.Value, flagRealm.Value)
			cmd.PrintErrln()
		},
	}
	parentCmd.AddCommand(cmd)
}
