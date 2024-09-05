package cmd

import (
	"os"
	"path/filepath"

	"github.com/riridotdev/sto-go/state"
	"github.com/riridotdev/sto-go/store"
	"github.com/spf13/cobra"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, dirName := filepath.Split(pwd)

	InitCmd.Flags().StringVar(&initName, "name", dirName, "name for new profile")
}

var initName string

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new Sto profile in the current directory",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		pwd, err := os.Getwd()
		if err != nil {
			panic("error getting pwd")
		}

		s, err := store.Init(pwd)
		if err != nil {
			return err
		}

		profile := state.Profile{
			Name: initName,
			Root: s.Root,
		}

		if err := State.AddProfile(profile); err != nil {
			return err
		}

		return nil
	},
}
