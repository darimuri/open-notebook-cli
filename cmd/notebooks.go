package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
	"github.com/darimuri/open-notebook-cli/internal/auth"
)

var (
	deleteSources bool
)

var notebooksCmd = &cobra.Command{
	Use:   "notebooks",
	Short: "Manage notebooks",
	Long:  `List, get, create, update, and delete notebooks`,
}

var notebooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notebooks",
	RunE:  runNotebooksList,
}

var notebooksGetCmd = &cobra.Command{
	Use:   "get [notebook_id]",
	Short: "Get a notebook by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksGet,
}

var notebooksCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new notebook",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksCreate,
}

var notebooksUpdateCmd = &cobra.Command{
	Use:   "update [notebook_id]",
	Short: "Update a notebook",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksUpdate,
}

var notebooksDeleteCmd = &cobra.Command{
	Use:   "delete [notebook_id]",
	Short: "Delete a notebook",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksDelete,
}

var notebooksDeletePreviewCmd = &cobra.Command{
	Use:   "delete-preview [notebook_id]",
	Short: "Preview what will be deleted",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksDeletePreview,
}

var notebooksAddSourceCmd = &cobra.Command{
	Use:   "add-source [notebook_id] [source_id]",
	Short: "Add a source to a notebook",
	Args:  cobra.ExactArgs(2),
	RunE:  runNotebooksAddSource,
}

var notebooksRemoveSourceCmd = &cobra.Command{
	Use:   "remove-source [notebook_id] [source_id]",
	Short: "Remove a source from a notebook",
	Args:  cobra.ExactArgs(2),
	RunE:  runNotebooksRemoveSource,
}

func init() {
	notebooksCmd.AddCommand(notebooksListCmd)
	notebooksCmd.AddCommand(notebooksGetCmd)
	notebooksCmd.AddCommand(notebooksCreateCmd)
	notebooksCmd.AddCommand(notebooksUpdateCmd)
	notebooksCmd.AddCommand(notebooksDeleteCmd)
	notebooksCmd.AddCommand(notebooksDeletePreviewCmd)
	notebooksCmd.AddCommand(notebooksAddSourceCmd)
	notebooksCmd.AddCommand(notebooksRemoveSourceCmd)
	rootCmd.AddCommand(notebooksCmd)

	notebooksDeleteCmd.Flags().BoolVar(&deleteSources, "delete-sources", false, "Delete exclusive sources when deleting notebook")
}

func getClient() *api.Client {
	cfg, _ := loadConfig()
	authMiddleware := auth.NewMiddleware(cfg.APIKey)
	return api.NewClientWithDebug(cfg.APIURL, authMiddleware, debug)
}

func runNotebooksList(cmd *cobra.Command, args []string) error {
	client := getClient()

	var notebooks []api.NotebookResponse
	err := client.Get("/api/notebooks", &notebooks)
	if err != nil {
		return fmt.Errorf("failed to list notebooks: %w", err)
	}

	return outputJSON(notebooks)
}

func runNotebooksGet(cmd *cobra.Command, args []string) error {
	client := getClient()

	var notebook api.NotebookResponse
	err := client.Get("/api/notebooks/"+args[0], &notebook)
	if err != nil {
		return fmt.Errorf("failed to get notebook: %w", err)
	}

	return outputJSON(notebook)
}

func runNotebooksCreate(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.NotebookCreate{Name: args[0]}
	var notebook api.NotebookResponse
	err := client.Post("/api/notebooks", req, &notebook)
	if err != nil {
		return fmt.Errorf("failed to create notebook: %w", err)
	}

	return outputJSON(notebook)
}

func runNotebooksUpdate(cmd *cobra.Command, args []string) error {
	client := getClient()

	if len(args) < 2 {
		return fmt.Errorf("usage: notebooks update [notebook_id] [new_name]")
	}

	name := args[1]
	req := api.NotebookUpdate{Name: &name}
	var notebook api.NotebookResponse
	err := client.Put("/api/notebooks/"+args[0], req, &notebook)
	if err != nil {
		return fmt.Errorf("failed to update notebook: %w", err)
	}

	return outputJSON(notebook)
}

func runNotebooksDelete(cmd *cobra.Command, args []string) error {
	client := getClient()

	var result api.NotebookDeleteResponse
	path := "/api/notebooks/" + args[0]
	if deleteSources {
		path += "?delete_exclusive_sources=true"
	}
	err := client.Delete(path, &result)
	if err != nil {
		return fmt.Errorf("failed to delete notebook: %w", err)
	}

	return outputJSON(result)
}

func runNotebooksDeletePreview(cmd *cobra.Command, args []string) error {
	client := getClient()

	var preview api.NotebookDeletePreview
	err := client.Get("/api/notebooks/"+args[0]+"/delete-preview", &preview)
	if err != nil {
		return fmt.Errorf("failed to get delete preview: %w", err)
	}

	return outputJSON(preview)
}

func runNotebooksAddSource(cmd *cobra.Command, args []string) error {
	client := getClient()

	notebookID := args[0]
	sourceID := args[1]
	var result any
	err := client.Post("/api/notebooks/"+notebookID+"/sources/"+sourceID, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to add source to notebook: %w", err)
	}

	return outputJSON(result)
}

func runNotebooksRemoveSource(cmd *cobra.Command, args []string) error {
	client := getClient()

	notebookID := args[0]
	sourceID := args[1]
	var result any
	err := client.Delete("/api/notebooks/"+notebookID+"/sources/"+sourceID, &result)
	if err != nil {
		return fmt.Errorf("failed to remove source from notebook: %w", err)
	}

	return outputJSON(result)
}