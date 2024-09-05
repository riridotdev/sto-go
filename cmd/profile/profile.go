package profile

import "github.com/spf13/cobra"

func init() {
	ProfileAddCmd.Flags().StringVar(&profileName, "name", "", "name for profile")

	ProfileCmd.AddCommand(
		ProfileAddCmd,
		ProfileSwitchCmd,
		ProfileListCmd,
		ProfileEnableCmd,
		ProfileDisableCmd,
		ProfileToggleCmd,
		ProfileRenameCmd,
	)
}

var ProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "",
	Args:  cobra.NoArgs,
}

var profileName string
