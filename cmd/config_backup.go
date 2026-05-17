package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration backup and restore",
	Long:  `Backup and restore models and credentials`,
}

var configBackupCmd = &cobra.Command{
	Use:   "backup [output_dir]",
	Short: "Backup models, credentials, and defaults",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigBackup,
}

var configRestoreCmd = &cobra.Command{
	Use:   "restore [input_dir]",
	Short: "Restore models, credentials, and defaults",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigRestore,
}

func init() {
	configCmd.AddCommand(configBackupCmd)
	configCmd.AddCommand(configRestoreCmd)
	rootCmd.AddCommand(configCmd)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func runConfigBackup(cmd *cobra.Command, args []string) error {
	client := getClient()
	outputDir := args[0]

	// Create output directory if not exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Backup credentials
	var credentials []api.CredentialResponse
	if err := client.Get("/api/credentials", &credentials); err != nil {
		return fmt.Errorf("failed to backup credentials: %w", err)
	}
	if err := writeJSON(outputDir+"/credentials.json", credentials); err != nil {
		return err
	}
	fmt.Printf("Backed up %d credentials\n", len(credentials))

	// Backup models
	var models []api.ModelResponse
	if err := client.Get("/api/models", &models); err != nil {
		return fmt.Errorf("failed to backup models: %w", err)
	}
	if err := writeJSON(outputDir+"/models.json", models); err != nil {
		return err
	}
	fmt.Printf("Backed up %d models\n", len(models))

	// Backup model defaults
	var defaults api.DefaultModelsResponse
	if err := client.Get("/api/models/defaults", &defaults); err != nil {
		return fmt.Errorf("failed to backup model defaults: %w", err)
	}
	if err := writeJSON(outputDir+"/models_defaults.json", defaults); err != nil {
		return err
	}
	fmt.Println("Backed up model defaults")

	fmt.Printf("\nBackup saved to: %s/\n", outputDir)
	return nil
}

func runConfigRestore(cmd *cobra.Command, args []string) error {
	client := getClient()
	inputDir := args[0]

	// Restore credentials
	var credentials []api.CredentialResponse
	if err := readJSON(inputDir+"/credentials.json", &credentials); err != nil {
		return fmt.Errorf("failed to read credentials backup: %w", err)
	}
	restoredCreds := 0
	for _, cred := range credentials {
		var result api.CredentialResponse
		err := client.Post("/api/credentials", cred, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore credential %s: %v\n", cred.Name, err)
		} else {
			restoredCreds++
		}
	}
	fmt.Printf("Restored %d credentials\n", restoredCreds)

	// Restore models
	var models []api.ModelResponse
	if err := readJSON(inputDir+"/models.json", &models); err != nil {
		return fmt.Errorf("failed to read models backup: %w", err)
	}
	restoredModels := 0
	for _, model := range models {
		var result api.ModelResponse
		err := client.Post("/api/models", model, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore model %s: %v\n", model.Name, err)
		} else {
			restoredModels++
		}
	}
	fmt.Printf("Restored %d models\n", restoredModels)

	// Restore model defaults
	var defaults api.DefaultModelsResponse
	if err := readJSON(inputDir+"/models_defaults.json", &defaults); err != nil {
		return fmt.Errorf("failed to read model defaults backup: %w", err)
	}
	if err := client.Put("/api/models/defaults", defaults, &defaults); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to restore model defaults: %v\n", err)
	} else {
		fmt.Println("Restored model defaults")
	}

	fmt.Println("\nRestore complete")
	return nil
}