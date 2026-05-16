package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
	"github.com/darimuri/open-notebook-cli/internal/auth"
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

func init() {
	notebooksCmd.AddCommand(notebooksListCmd)
	notebooksCmd.AddCommand(notebooksGetCmd)
	notebooksCmd.AddCommand(notebooksCreateCmd)
	notebooksCmd.AddCommand(notebooksUpdateCmd)
	notebooksCmd.AddCommand(notebooksDeleteCmd)
	notebooksCmd.AddCommand(notebooksDeletePreviewCmd)
	rootCmd.AddCommand(notebooksCmd)
}

func getClient() *api.Client {
	cfg, _ := loadConfig()
	authMiddleware := auth.NewMiddleware(cfg.APIKey)
	return api.NewClient(cfg.APIURL, authMiddleware)
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

	// For now, just update name from args
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
	err := client.Delete("/api/notebooks/"+args[0], &result)
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

