package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Manage sources",
	Long:  `List, upload, download, and manage sources`,
}

var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sources",
	RunE:  runSourcesList,
}

var sourcesUploadCmd = &cobra.Command{
	Use:   "upload [file_path]",
	Short: "Upload a source file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSourcesUpload,
}

var sourcesDownloadCmd = &cobra.Command{
	Use:   "download [source_id]",
	Short: "Download a source",
	Args:  cobra.ExactArgs(1),
	RunE:  runSourcesDownload,
}

var sourcesRetryCmd = &cobra.Command{
	Use:   "retry [source_id]",
	Short: "Retry a failed source",
	Args:  cobra.ExactArgs(1),
	RunE:  runSourcesRetry,
}

func init() {
	sourcesCmd.AddCommand(sourcesListCmd)
	sourcesCmd.AddCommand(sourcesUploadCmd)
	sourcesCmd.AddCommand(sourcesDownloadCmd)
	sourcesCmd.AddCommand(sourcesRetryCmd)
	rootCmd.AddCommand(sourcesCmd)
}

func runSourcesList(cmd *cobra.Command, args []string) error {
	client := getClient()

	var sources []api.SourceResponse
	err := client.Get("/api/sources", &sources)
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	return outputJSON(sources)
}

func runSourcesUpload(cmd *cobra.Command, args []string) error {
	client := getClient()

	// Read file content
	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// For now, just send as JSON - real implementation would use multipart
	req := map[string]interface{}{
		"filename": args[0],
		"content":  string(content),
	}
	var result api.SourceResponse
	err = client.Post("/api/sources/json", req, &result)
	if err != nil {
		return fmt.Errorf("failed to upload source: %w", err)
	}

	return outputJSON(result)
}

func runSourcesDownload(cmd *cobra.Command, args []string) error {
	client := getClient()

	var source api.SourceResponse
	err := client.Get("/api/sources/"+args[0]+"/download", &source)
	if err != nil {
		return fmt.Errorf("failed to download source: %w", err)
	}

	return outputJSON(source)
}

func runSourcesRetry(cmd *cobra.Command, args []string) error {
	client := getClient()

	var result api.SourceResponse
	err := client.Post("/api/sources/"+args[0]+"/retry", nil, &result)
	if err != nil {
		return fmt.Errorf("failed to retry source: %w", err)
	}

	return outputJSON(result)
}

