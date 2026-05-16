package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search and ask",
	Long:  `Search notebooks, ask questions`,
}

var searchSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search notebooks",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchSearch,
}

var searchAskCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchAsk,
}

var searchSimpleAskCmd = &cobra.Command{
	Use:   "simple [question]",
	Short: "Simple ask (quick answer)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchSimpleAsk,
}

func init() {
	searchCmd.AddCommand(searchSearchCmd)
	searchCmd.AddCommand(searchAskCmd)
	searchCmd.AddCommand(searchSimpleAskCmd)
	rootCmd.AddCommand(searchCmd)
}

func runSearchSearch(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.SearchRequest{Query: args[0]}
	var results api.SearchResponse
	err := client.Post("/api/search", req, &results)
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	return outputJSON(results)
}

func runSearchAsk(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.AskRequest{Question: args[0]}
	var answer api.AskResponse
	err := client.Post("/api/search/ask", req, &answer)
	if err != nil {
		return fmt.Errorf("failed to ask: %w", err)
	}

	return outputJSON(answer)
}

func runSearchSimpleAsk(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.AskRequest{Question: args[0]}
	var answer api.AskResponse
	err := client.Post("/api/search/ask/simple", req, &answer)
	if err != nil {
		return fmt.Errorf("failed to ask: %w", err)
	}

	return outputJSON(answer)
}

