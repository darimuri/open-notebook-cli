package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
)

var searchNotebook string
var strategyModel string
var answerModel string
var finalAnswerModel string

var defaultModels struct {
	strategyModel     string
	answerModel       string
	finalAnswerModel  string
}
var modelsLoaded bool

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
	Short: "Ask a question (simple, non-streaming)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchAsk,
}

var searchSimpleAskCmd = &cobra.Command{
	Use:   "simple [question]",
	Short: "Simple ask (non-streaming)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearchSimpleAsk,
}

func init() {
	searchCmd.AddCommand(searchSearchCmd)
	searchCmd.AddCommand(searchAskCmd)
	searchCmd.AddCommand(searchSimpleAskCmd)
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchNotebook, "notebook", "n", "", "Notebook ID to search in")
	searchCmd.Flags().StringVar(&strategyModel, "strategy-model", "", "Model ID for query strategy")
	searchCmd.Flags().StringVar(&answerModel, "answer-model", "", "Model ID for individual answers")
	searchCmd.Flags().StringVar(&finalAnswerModel, "final-answer-model", "", "Model ID for final answer")
}

func loadDefaultModels() error {
	if modelsLoaded {
		return nil
	}
	client := getClient()
	var models []api.ModelResponse
	if err := client.Get("/api/models", &models); err != nil {
		return err
	}
	for _, m := range models {
		if m.Type == "language" && defaultModels.strategyModel == "" {
			defaultModels.strategyModel = m.ID
			defaultModels.answerModel = m.ID
			defaultModels.finalAnswerModel = m.ID
			break
		}
	}
	modelsLoaded = true
	return nil
}

func runSearchSearch(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.SearchRequest{Query: args[0]}
	notebookID := searchNotebook
	if notebookID == "" {
		notebookID = getDefaultNotebook()
	}
	if notebookID != "" {
		req.NotebookIDs = []string{notebookID}
	}
	var results api.SearchResponse
	err := client.Post("/api/search", req, &results)
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	return outputJSON(results)
}

func runSearchAsk(cmd *cobra.Command, args []string) error {
	if err := loadDefaultModels(); err != nil {
		return fmt.Errorf("failed to load default models: %w", err)
	}

	client := getClient()

	req := api.AskRequest{Question: args[0]}
	if strategyModel != "" {
		req.StrategyModel = strategyModel
	} else {
		req.StrategyModel = defaultModels.strategyModel
	}
	if answerModel != "" {
		req.AnswerModel = answerModel
	} else {
		req.AnswerModel = defaultModels.answerModel
	}
	if finalAnswerModel != "" {
		req.FinalAnswerModel = finalAnswerModel
	} else {
		req.FinalAnswerModel = defaultModels.finalAnswerModel
	}
	notebookID := searchNotebook
	if notebookID == "" {
		notebookID = getDefaultNotebook()
	}
	if notebookID != "" {
		req.NotebookIDs = []string{notebookID}
	}
	var answer api.AskResponse
	err := client.Post("/api/search/ask/simple", req, &answer)
	if err != nil {
		return fmt.Errorf("failed to ask: %w", err)
	}

	return outputJSON(answer)
}

func runSearchSimpleAsk(cmd *cobra.Command, args []string) error {
	if err := loadDefaultModels(); err != nil {
		return fmt.Errorf("failed to load default models: %w", err)
	}

	client := getClient()

	req := api.AskRequest{Question: args[0]}
	if strategyModel != "" {
		req.StrategyModel = strategyModel
	} else {
		req.StrategyModel = defaultModels.strategyModel
	}
	if answerModel != "" {
		req.AnswerModel = answerModel
	} else {
		req.AnswerModel = defaultModels.answerModel
	}
	if finalAnswerModel != "" {
		req.FinalAnswerModel = finalAnswerModel
	} else {
		req.FinalAnswerModel = defaultModels.finalAnswerModel
	}
	notebookID := searchNotebook
	if notebookID == "" {
		notebookID = getDefaultNotebook()
	}
	if notebookID != "" {
		req.NotebookIDs = []string{notebookID}
	}
	var simpleAnswer api.AskResponse
	err := client.Post("/api/search/ask/simple", req, &simpleAnswer)
	if err != nil {
		return fmt.Errorf("failed to ask: %w", err)
	}

	return outputJSON(simpleAnswer)
}

