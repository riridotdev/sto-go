package profile

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/riridotdev/sto-go/cmd"
)

var ProfileListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Args:  cobra.NoArgs,
	RunE: func(c *cobra.Command, args []string) error {
		writer := tabwriter.NewWriter(c.OutOrStdout(), 1, 1, 1, ' ', 0)

		writer.Write([]byte("STATUS\tNAME\tLOCATION\n"))

		for _, profile := range cmd.State.Profiles {
			if profile.Name == cmd.State.DefaultProfile {
				writer.Write([]byte("[Default]\t"))
			} else if profile.Active {
				writer.Write([]byte("[Active]\t"))
			} else {
				writer.Write([]byte("[Inactive]\t"))
			}
			writer.Write([]byte(fmt.Sprintf("%s\t(%s)\n", profile.Name, profile.Root)))
		}

		return writer.Flush()
	},
}
