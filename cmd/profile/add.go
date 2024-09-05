package profile

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
	"github.com/riridotdev/sto-go/state"
	"github.com/riridotdev/sto-go/store"
)

var ProfileAddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		targetPath := args[0]

		targetPath, err := homedir.Expand(targetPath)
		if err != nil {
			return err
		}
		targetPath, err = filepath.Abs(targetPath)
		if err != nil {
			return err
		}

		if profileName == "" {
			_, profileName = filepath.Split(targetPath)
		}

		s, err := store.Restore(targetPath)
		if err != nil {
			return err
		}

		profile := state.Profile{
			Name: profileName,
			Root: s.Root,
		}

		if err := cmd.State.AddProfile(profile); err != nil {
			return err
		}

		fmt.Printf("Added profile: %s\n", profileName)

		return nil
	},
}
