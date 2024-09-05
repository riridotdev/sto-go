package profile

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		profileName := args[0]

		profile, err := cmd.State.GetProfile(profileName)
		if err != nil {
			return err
		}

		profile.Active = !profile.Active

		if profile.Active {
			fmt.Printf("Enabled profile %q\n", profileName)
		} else {
			fmt.Printf("Disabled profile %q\n", profileName)
		}

		return nil
	},
}
