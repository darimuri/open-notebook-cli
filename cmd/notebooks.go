package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

var notebooksDeleteCmd = &cobra.Command{
	Use:   "delete [notebook_id]",
	Short: "Delete a notebook",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotebooksDelete,
}

func init() {
	notebooksCmd.AddCommand(notebooksListCmd)
	notebooksCmd.AddCommand(notebooksGetCmd)
	notebooksCmd.AddCommand(notebooksCreateCmd)
	notebooksCmd.AddCommand(notebooksDeleteCmd)
	rootCmd.AddCommand(notebooksCmd)
}

func runNotebooksList(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notebooks list)\n", cfg.APIURL)
	return nil
}

func runNotebooksGet(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notebooks get: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runNotebooksCreate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notebooks create: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runNotebooksDelete(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notebooks delete: %s)\n", cfg.APIURL, args[0])
	return nil
}