package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (search: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runSearchAsk(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (ask: %s)\n", cfg.APIURL, args[0])
	return nil
}

func runSearchSimpleAsk(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	fmt.Printf("API: %s (simple-ask: %s)\n", cfg.APIURL, args[0])
	return nil
}