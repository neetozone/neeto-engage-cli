package commands

import (
	"github.com/spf13/cobra"
)

var changelogsCmd = &cobra.Command{
	Use:   "changelogs",
	Short: "Show changelogs",
}

var changelogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List changelogs that are not archived, newest first",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := paginationParams(cmd)
		if query, _ := cmd.Flags().GetString("query"); query != "" {
			params.Set("query", query)
		}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/changelogs", params)
		if err != nil {
			return err
		}

		printList(data, "changelogs", nil)
		return nil
	},
}

func init() {
	changelogsListCmd.Flags().String("query", "", "Filter by title (partial match)")
	addPaginationFlags(changelogsListCmd)

	changelogsCmd.AddCommand(changelogsListCmd)
	register(func(root *cobra.Command) { root.AddCommand(changelogsCmd) })
}
