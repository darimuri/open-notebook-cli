package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
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
	Use:   "create [notebook_id] [content]",
	Short: "Create a new note",
	Args:  cobra.ExactArgs(2),
	RunE:  runNotesCreate,
}

var notesUpdateCmd = &cobra.Command{
	Use:   "update [note_id]",
	Short: "Update a note",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotesUpdate,
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
	notesCmd.AddCommand(notesUpdateCmd)
	notesCmd.AddCommand(notesDeleteCmd)
	rootCmd.AddCommand(notesCmd)
}

func runNotesList(cmd *cobra.Command, args []string) error {
	client := getClient()

	var notes []api.NoteResponse
	err := client.Get("/api/notes", &notes)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return outputJSON(notes)
}

func runNotesGet(cmd *cobra.Command, args []string) error {
	client := getClient()

	var note api.NoteResponse
	err := client.Get("/api/notes/"+args[0], &note)
	if err != nil {
		return fmt.Errorf("failed to get note: %w", err)
	}

	return outputJSON(note)
}

func runNotesCreate(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.NoteCreate{
		NotebookID: args[0],
		Content:    args[1],
	}
	var note api.NoteResponse
	err := client.Post("/api/notes", req, &note)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return outputJSON(note)
}

func runNotesUpdate(cmd *cobra.Command, args []string) error {
	client := getClient()

	// For now, just update content from args
	content := args[1]
	req := api.NoteUpdate{Content: &content}
	var note api.NoteResponse
	err := client.Put("/api/notes/"+args[0], req, &note)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return outputJSON(note)
}

func runNotesDelete(cmd *cobra.Command, args []string) error {
	client := getClient()

	err := client.Delete("/api/notes/"+args[0], nil)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Println("Note deleted successfully")
	return nil
}

