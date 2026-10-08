package commands

import (
	"github.com/spf13/cobra"
)

var votesCmd = &cobra.Command{
	Use:   "votes",
	Short: "Show the feature requests a customer voted for",
}

var votesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the feature requests a customer voted for",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		params := paginationParams(cmd)
		params.Set("email", email)

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/votes", params)
		if err != nil {
			return err
		}

		printList(data, "feature_requests", nil)
		return nil
	},
}

func init() {
	votesListCmd.Flags().String("email", "", "Email address of the customer")
	_ = votesListCmd.MarkFlagRequired("email")
	addPaginationFlags(votesListCmd)

	votesCmd.AddCommand(votesListCmd)
	register(func(root *cobra.Command) { root.AddCommand(votesCmd) })
}
