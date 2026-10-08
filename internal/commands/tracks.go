package commands

import (
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "Show the columns of the feature request board",
}

var tracksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tracks, the columns of the feature request board",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/tracks", nil)
		if err != nil {
			return err
		}

		printList(data, "tracks", []output.Breadcrumb{
			{Label: "Move", Command: "neetoengage feature-requests move <id> --track-id <track-id>"},
		})
		return nil
	},
}

func init() {
	tracksCmd.AddCommand(tracksListCmd)
	register(func(root *cobra.Command) { root.AddCommand(tracksCmd) })
}
