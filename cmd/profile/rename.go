package profile

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		profileName := args[0]
		newName := args[1]

		profile, err := cmd.State.GetProfile(profileName)
		if err != nil {
			return err
		}

		profile.Name = newName

		fmt.Printf("Renamed profile %q to %q\n", profileName, newName)

		return nil
	},
}
