package profile

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		profileName := args[0]

		profile, err := cmd.State.GetProfile(profileName)
		if err != nil {
			return err
		}

		if err := profile.Enable(); err != nil {
			return err
		}

		fmt.Printf("Enabled profile %q\n", profileName)

		return nil
	},
}
