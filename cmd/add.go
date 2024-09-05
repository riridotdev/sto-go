package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/mitchellh/go-homedir"
	"github.com/riridotdev/sto-go/store"
	"github.com/spf13/cobra"
)

func init() {
	AddCmd.Flags().StringVar(&addDestination, "destination", "", "destination for the source file to be linked to")
}

var addDestination string

var AddCmd = &cobra.Command{
	Use:   "add",
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

		targetPath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		_, entryName := filepath.Split(targetPath)
		entrySourcePath := fmt.Sprintf("%s/%s", s.Root, entryName)

		newEntry := store.Entry{
			Name:        entryName,
			Source:      entrySourcePath,
			Destination: targetPath,
		}

		if addDestination != "" {
			if addDestination, err = homedir.Expand(addDestination); err != nil {
				return err
			}
			if addDestination, err = filepath.Abs(addDestination); err != nil {
				return err
			}
			newEntry.Destination = addDestination
		}

		if err = s.Add(newEntry); err != nil {
			return err
		}

		if err := os.Rename(targetPath, entrySourcePath); err != nil {
			return err
		}

		if profile.Active {
			if err := newEntry.Link(); err != nil {
				os.Rename(newEntry.Source, targetPath)
				return err
			}
		}

		writer := tabwriter.NewWriter(cmd.OutOrStdout(), 1, 1, 1, ' ', 0)
		writer.Write([]byte(fmt.Sprintf("[Linked]\t%s\t%s\t->\t%s\n", newEntry.Name, newEntry.Source, newEntry.Destination)))

		if err := writer.Flush(); err != nil {
			return err
		}

		if err := s.Persist(); err != nil {
			return err
		}

		return nil
	},
}
