package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/state"
	"github.com/riridotdev/sto-go/store"
)

func init() {
	ListCmd.Flags().BoolVarP(&listAll, "all", "A", false, "show results for all active profiles")
}

var listAll bool

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Profile == "" {
			return ErrProfileNotSet
		}

		var profiles []state.Profile
		if listAll {
			profiles = State.Profiles
		} else {
			profiles = State.ActiveProfiles()
		}

		writer := tabwriter.NewWriter(cmd.OutOrStdout(), 1, 1, 1, ' ', 0)

		for i, profile := range profiles {
			s, err := store.Restore(profile.Root)
			if err != nil {
				return err
			}

			entries := s.Entries()

			if len(entries) == 0 {
				continue
			}

			if len(profiles) > 1 {
				fmt.Printf("Profile: %s\n", profile.Name)
			}

			for _, entry := range entries {
				entryState, err := entry.Check()
				if err != nil {
					// TODO: Handle errors
					continue
				}

				switch entryState {
				case store.Linked:
					writer.Write([]byte(fmt.Sprintf("[Linked]\t%s:\t%s\t->\t%s\n", entry.Name, entry.Source, entry.Destination)))
				case store.Unlinked:
					writer.Write([]byte(fmt.Sprintf("[Unlinked]\t%s:\t%s\t->\t%s\n", entry.Name, entry.Source, entry.Destination)))
				default:
					writer.Write([]byte(fmt.Sprintf("[Unknown]\t%s:\t%s\t->\t%s\n", entry.Name, entry.Source, entry.Destination)))
				}
			}

			if err := writer.Flush(); err != nil {
				return err
			}

			if i != len(profiles)-1 {
				print("\n")
			}
		}

		return nil
	},
}
