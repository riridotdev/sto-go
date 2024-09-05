package profile

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		profileName := args[0]

		profile, err := cmd.State.GetProfile(profileName)
		if err != nil {
			return err
		}

		if err := profile.Disable(); err != nil {
			return err
		}

		fmt.Printf("Disabled profile %q\n", profileName)

		return nil
	},
}
