package commands

import (
	"net/url"
	"strconv"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var featureRequestsCmd = &cobra.Command{
	Use:   "feature-requests",
	Short: "Search, file and follow up on feature requests",
}

var featureRequestsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search public feature requests by title, description and tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		params := url.Values{"query": {query}}
		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			params.Set("per_page", strconv.Itoa(limit))
		}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/feature_requests", params)
		if err != nil {
			return err
		}

		printList(data, "feature_requests", []output.Breadcrumb{
			{Label: "Add voter", Command: "neetoengage feature-requests add-voter <id> --customer-email <email> --note <note>"},
			{Label: "Create", Command: "neetoengage feature-requests create --title <title> --description <description> --customer-email <email> --note <note>"},
		})
		return nil
	},
}

var featureRequestsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a feature request for a customer, with the customer as its first voter and a private note",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		body := customerPayload(cmd)
		body["feature_request"] = map[string]interface{}{"title": title, "description": description}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post("/feature_requests", body)
		if err != nil {
			return err
		}

		printResource(data, []output.Breadcrumb{
			{Label: "Voters", Command: "neetoengage feature-requests voters <id>"},
		})
		return nil
	},
}

var featureRequestsAddVoterCmd = &cobra.Command{
	Use:   "add-voter <id>",
	Short: "Add a customer as a voter on a feature request and attach a private note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post(featureRequestPath(args[0], "upvotes"), customerPayload(cmd))
		if err != nil {
			return err
		}

		printResource(data, []output.Breadcrumb{
			{Label: "Voters", Command: "neetoengage feature-requests voters " + args[0]},
		})
		return nil
	},
}

var featureRequestsVotersCmd = &cobra.Command{
	Use:   "voters <id>",
	Short: "List the customers who voted for a feature request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get(featureRequestPath(args[0], "voters"), paginationParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "voters", nil)
		return nil
	},
}

var featureRequestsMoveCmd = &cobra.Command{
	Use:   "move <id>",
	Short: "Move a feature request to another track, optionally attaching a changelog and emailing voters",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trackID, _ := cmd.Flags().GetString("track-id")
		notify, _ := cmd.Flags().GetBool("notify-all-voters")
		body := map[string]interface{}{"track_id": trackID, "notify_all_voters": notify}
		if v, _ := cmd.Flags().GetString("message"); v != "" {
			body["message"] = v
		}
		if v, _ := cmd.Flags().GetString("changelog-id"); v != "" {
			body["changelog_id"] = v
		}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Patch(featureRequestPath(args[0], "track"), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var featureRequestsCommentCmd = &cobra.Command{
	Use:   "comment <id>",
	Short: "Post a public comment on a feature request, optionally emailing it to voters",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, _ := cmd.Flags().GetString("content")
		notify, _ := cmd.Flags().GetBool("notify-all-voters")

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Post(featureRequestPath(args[0], "comments"), map[string]interface{}{
			"comment":           map[string]interface{}{"content": content},
			"notify_all_voters": notify,
		})
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

func featureRequestPath(id, resource string) string {
	return "/feature_requests/" + url.PathEscape(id) + "/" + resource
}

func customerPayload(cmd *cobra.Command) map[string]interface{} {
	email, _ := cmd.Flags().GetString("customer-email")
	name, _ := cmd.Flags().GetString("customer-name")
	note, _ := cmd.Flags().GetString("note")
	return map[string]interface{}{
		"customer":       map[string]interface{}{"email": email, "name": name},
		"reference_note": map[string]interface{}{"body": note},
	}
}

func addCustomerFlags(cmd *cobra.Command) {
	cmd.Flags().String("customer-email", "", "Email address of the customer who asked for it")
	cmd.Flags().String("customer-name", "", "Full name of the customer who asked for it")
	cmd.Flags().String("note", "", "Private note for the team. Customers never see it.")
	_ = cmd.MarkFlagRequired("customer-email")
	_ = cmd.MarkFlagRequired("note")
}

func init() {
	featureRequestsSearchCmd.Flags().String("query", "", "Words to search for")
	featureRequestsSearchCmd.Flags().Int("limit", 0, "Maximum requests to return (default 8, max 15)")
	_ = featureRequestsSearchCmd.MarkFlagRequired("query")

	featureRequestsCreateCmd.Flags().String("title", "", "Short title, at most 120 characters. Shown publicly.")
	featureRequestsCreateCmd.Flags().String("description", "", "Short description of the request. Shown publicly.")
	_ = featureRequestsCreateCmd.MarkFlagRequired("title")
	_ = featureRequestsCreateCmd.MarkFlagRequired("description")
	addCustomerFlags(featureRequestsCreateCmd)

	addCustomerFlags(featureRequestsAddVoterCmd)

	addPaginationFlags(featureRequestsVotersCmd)

	featureRequestsMoveCmd.Flags().String("track-id", "", "ID of the track to move the feature request to")
	featureRequestsMoveCmd.Flags().Bool("notify-all-voters", false, "Email the status update to every voter")
	featureRequestsMoveCmd.Flags().String("message", "", "HTML body of the status update email. Defaults to the \"now live\" announcement.")
	featureRequestsMoveCmd.Flags().String("changelog-id", "", "ID of a changelog to attach to the feature request")
	_ = featureRequestsMoveCmd.MarkFlagRequired("track-id")

	featureRequestsCommentCmd.Flags().String("content", "", "Text of the comment. Shown publicly.")
	featureRequestsCommentCmd.Flags().Bool("notify-all-voters", false, "Email the comment to every voter")
	_ = featureRequestsCommentCmd.MarkFlagRequired("content")

	featureRequestsCmd.AddCommand(featureRequestsSearchCmd)
	featureRequestsCmd.AddCommand(featureRequestsCreateCmd)
	featureRequestsCmd.AddCommand(featureRequestsAddVoterCmd)
	featureRequestsCmd.AddCommand(featureRequestsVotersCmd)
	featureRequestsCmd.AddCommand(featureRequestsMoveCmd)
	featureRequestsCmd.AddCommand(featureRequestsCommentCmd)
	register(func(root *cobra.Command) { root.AddCommand(featureRequestsCmd) })
}
