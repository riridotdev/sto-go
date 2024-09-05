package profile

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		targetName := args[0]

		if cmd.State.DefaultProfile == targetName {
			fmt.Printf("Default profile already set to %q\n", targetName)
			return nil
		}

		newProfile, err := cmd.State.GetProfile(targetName)
		if err != nil {
			return err
		}

		activeProfiles := cmd.State.ActiveProfiles()

		if len(activeProfiles) == 1 {
			oldProfile, err := cmd.State.GetProfile(cmd.State.DefaultProfile)
			if err != nil {
				return err
			}

			if err := oldProfile.Disable(); err != nil {
				return err
			}

			if err := newProfile.Enable(); err != nil {
				return err
			}
		}

		if !newProfile.Active {
			return fmt.Errorf("profile %q must be enabled before switching", targetName)
		}

		cmd.State.DefaultProfile = newProfile.Name

		fmt.Printf("Default profile set: %s\n", newProfile.Name)

		return nil
	},
}
