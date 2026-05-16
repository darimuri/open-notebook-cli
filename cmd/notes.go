package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Manage notes",
	Long:  `List, get, create, update, and delete notes`,
}

var notesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	RunE:  runNotesList,
}

var notesGetCmd = &cobra.Command{
	Use:   "get [note_id]",
	Short: "Get a note by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotesGet,
}

var notesCreateCmd = &cobra.Command{
	Use:   "create [content]",
	Short: "Create a new note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotesCreate,
}

var notesDeleteCmd = &cobra.Command{
	Use:   "delete [note_id]",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotesDelete,
}

func init() {
	notesCmd.AddCommand(notesListCmd)
	notesCmd.AddCommand(notesGetCmd)
	notesCmd.AddCommand(notesCreateCmd)
	notesCmd.AddCommand(notesDeleteCmd)
	rootCmd.AddCommand(notesCmd)
}

func runNotesList(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notes list)\n", cfg.APIURL)
	return nil
}

func runNotesGet(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notes get: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runNotesCreate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notes create: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runNotesDelete(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (notes delete: %s)\n", cfg.APIURL, args[0])
	return nil
}