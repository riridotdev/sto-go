package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/riridotdev/sto-go/store"
	"github.com/spf13/cobra"
)

var UnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if Profile == "" {
			return ErrProfileNotSet
		}

		profile, err := State.GetProfile(Profile)
		if err != nil {
			return err
		}

		s, err := store.Restore(profile.Root)
		if err != nil {
			return err
		}

		targetEntry := args[0]

		entry, err := s.Entry(targetEntry)
		if err != nil {
			return err
		}

		entryState, err := entry.Check()
		if err != nil {
			return err
		}

		if entryState != store.Linked {
			_ = profile.RemoveLinkedEntry(profile.Name)
			fmt.Printf("Entry %q already unlinked\n", entry.Name)
			return nil
		}

		if err := entry.Unlink(); err != nil {
			return err
		}

		if err := profile.RemoveLinkedEntry(profile.Name); err != nil {
			return err
		}

		writer := tabwriter.NewWriter(cmd.OutOrStdout(), 1, 1, 1, ' ', 0)
		writer.Write([]byte(fmt.Sprintf("[Unlinked]\t%s:\t%s\t->\t%s\n", entry.Name, entry.Source, entry.Destination)))
		writer.Flush()

		return nil
	},
}
