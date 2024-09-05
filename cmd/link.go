package cmd

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/riridotdev/sto-go/state"
	"github.com/riridotdev/sto-go/store"
	"github.com/spf13/cobra"
)

var LinkCmd = &cobra.Command{
	Use:   "link",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if entryState != store.Linked && entryState != store.Unlinked {
			_ = profile.RemoveLinkedEntry(entry.Name)
			return fmt.Errorf("Error linking %q:\n%s", entry.Name, entryState.Describe(entry))
		}

		var errDuplicateEntry state.ErrDuplicateEntry
		if err := profile.AddLinkedEntry(entry.Name); err != nil && !errors.As(err, &errDuplicateEntry) {
			return err
		}

		if err := entry.Link(); err != nil {
			return err
		}

		writer := tabwriter.NewWriter(cmd.OutOrStdout(), 1, 1, 1, ' ', 0)
		writer.Write([]byte(fmt.Sprintf("[Linked]\t%s:\t%s\t->\t%s\n", entry.Name, entry.Source, entry.Destination)))

		return writer.Flush()
	},
}
