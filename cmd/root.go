package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/config"
)

var (
	cfgFile   string
	output    string
	apiURL    string
	apiKey    string
	notebook  string
)

var rootCmd = &cobra.Command{
	Use:   "open-notebook",
	Short: "CLI for Open Notebook API",
	Long:  `CLI tool for Open Notebook - Research Assistant API`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format (table, json)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API server URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key")
	rootCmd.PersistentFlags().StringVar(&notebook, "notebook", "", "default notebook ID")

	rootCmd.AddCommand(notebooksCmd)
	rootCmd.AddCommand(notesCmd)
	rootCmd.AddCommand(sourcesCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(skillsCmd)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	if apiURL != "" {
		cfg.APIURL = apiURL
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	if notebook != "" {
		cfg.Notebook = notebook
	}

	return cfg, nil
}

func getDefaultNotebook() string {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return ""
	}
	return cfg.Notebook
}