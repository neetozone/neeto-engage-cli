package commands

import (
	"errors"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show and update the product name and website URL",
}

var settingsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the product name and website URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Get("/setting", nil)
		if err != nil {
			return err
		}

		printResource(data, []output.Breadcrumb{
			{Label: "Update", Command: "neetoengage settings update --product-name <name> --website-url <url>"},
		})
		return nil
	},
}

var settingsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the product name, the website URL, or both",
	RunE: func(cmd *cobra.Command, args []string) error {
		setting := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("product-name"); cmd.Flags().Changed("product-name") {
			setting["product_name"] = v
		}
		if v, _ := cmd.Flags().GetString("website-url"); cmd.Flags().Changed("website-url") {
			setting["website_url"] = v
		}
		if len(setting) == 0 {
			return errors.New("pass --product-name, --website-url, or both")
		}

		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		data, err := c.Patch("/setting", map[string]interface{}{"setting": setting})
		if err != nil {
			return err
		}

		printResource(data, []output.Breadcrumb{
			{Label: "Show", Command: "neetoengage settings show"},
		})
		return nil
	},
}

func init() {
	settingsUpdateCmd.Flags().String("product-name", "", "New product name")
	settingsUpdateCmd.Flags().String("website-url", "", "New website URL, starting with http://, https:// or www.")

	settingsCmd.AddCommand(settingsShowCmd)
	settingsCmd.AddCommand(settingsUpdateCmd)
	register(func(root *cobra.Command) { root.AddCommand(settingsCmd) })
}
