package main

import (
	"errors"
	"fmt"

	"github.com/adrg/xdg"
	"github.com/riridotdev/sto-go/cmd"
	"github.com/riridotdev/sto-go/cmd/profile"
	"github.com/riridotdev/sto-go/state"
	"github.com/spf13/cobra"
)

func init() {
	err := readState()
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&cmd.Profile, "profile", cmd.State.DefaultProfile, "profile to use")
}

func readState() error {
	stateFilePath := fmt.Sprintf("%s/sto/.state", xdg.StateHome)

	st, err := state.Restore(stateFilePath)
	if err != nil {
		var errStatFileNotFound state.ErrStateFileNotFound
		if !errors.As(err, &errStatFileNotFound) {
			return err
		}
		st = state.State{StateFilePath: stateFilePath}
	}
	cmd.State = st

	return nil
}

var rootCmd = cobra.Command{}

func main() {
	rootCmd.AddCommand(
		cmd.InitCmd,
		cmd.ListCmd,
		cmd.AddCmd,
		cmd.LinkCmd,
		cmd.UnlinkCmd,
		profile.ProfileCmd,
	)
	rootCmd.SilenceUsage = true

	_ = rootCmd.Execute()

	if err := cmd.State.Persist(); err != nil {
		panic(err)
	}
}
