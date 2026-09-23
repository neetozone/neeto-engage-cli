package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neetozone/neeto-engage-cli/internal/auth"
	"github.com/neetozone/neeto-engage-cli/internal/client"
	"github.com/neetozone/neeto-engage-cli/internal/output"
	"github.com/spf13/cobra"
)

func getClient(cmd *cobra.Command) (*client.Client, error) {
	subdomain, _ := cmd.Flags().GetString("subdomain")
	creds, err := auth.SelectCredentials(subdomain)
	if err != nil {
		return nil, err
	}
	return client.New(creds), nil
}

func printResource(data json.RawMessage, breadcrumbs []output.Breadcrumb) {
	output.Print(data, breadcrumbs)
}

func readJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", path, err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return result, nil
}
