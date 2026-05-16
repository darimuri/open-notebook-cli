package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (sources list)\n", cfg.APIURL)
	return nil
}

func runSourcesUpload(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (sources upload: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runSourcesDownload(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (sources download: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runSourcesRetry(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (sources retry: %s)\n", cfg.APIURL, args[0])
	return nil
}