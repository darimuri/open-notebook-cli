package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	commandFilter string
	statusFilter string
	limitCount   int
)

var commandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "Manage and check command status",
}

var commandsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent commands",
	RunE:  runCommandsList,
}

var commandsStatusCmd = &cobra.Command{
	Use:   "status [command_id]",
	Short: "Get status of a command",
	Args:  cobra.ExactArgs(1),
	RunE:  runCommandsStatus,
}

func init() {
	commandsCmd.AddCommand(commandsListCmd)
	commandsCmd.AddCommand(commandsStatusCmd)
	rootCmd.AddCommand(commandsCmd)

	commandsListCmd.Flags().StringVar(&commandFilter, "command", "", "Filter by command name")
	commandsListCmd.Flags().StringVar(&statusFilter, "status", "", "Filter by status (new, running, completed, failed)")
	commandsListCmd.Flags().IntVar(&limitCount, "limit", 50, "Maximum number of jobs to return")
}

func runCommandsList(cmd *cobra.Command, args []string) error {
	client := getClient()

	// Build query string
	query := "/api/commands/jobs?"
	if commandFilter != "" {
		query += "command_filter=" + commandFilter + "&"
	}
	if statusFilter != "" {
		query += "status_filter=" + statusFilter + "&"
	}
	query += "limit=" + fmt.Sprintf("%d", limitCount)

	var result []map[string]interface{}
	err := client.Get(query, &result)
	if err != nil {
		return fmt.Errorf("failed to list commands: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no commands found (server v1.8.5 may not have list_command_jobs implemented)")
	}

	fmt.Printf("%-50s %-12s %s\n", "COMMAND_ID", "STATUS", "COMMAND_NAME")
	fmt.Println("------------------------------------------------------------------------")
	for _, cmd := range result {
		jobID := getStringValue(cmd, "job_id")
		status := getStringValue(cmd, "status")
		commandName := ""
		if result, ok := cmd["result"].(map[string]interface{}); ok {
			commandName = getStringValue(result, "command_name")
		}
		fmt.Printf("%-50s %-12s %s\n", jobID, status, commandName)
	}

	return nil
}

func runCommandsStatus(cmd *cobra.Command, args []string) error {
	client := getClient()

	commandID := args[0]

	var result map[string]interface{}
	err := client.Get("/api/commands/jobs/"+commandID, &result)
	if err != nil {
		return fmt.Errorf("failed to get command status: %w", err)
	}

	jobID := getStringValue(result, "job_id")
	status := getStringValue(result, "status")
	errorMsg := getStringValue(result, "error_message")

	fmt.Printf("Job ID: %s\n", jobID)
	fmt.Printf("Status: %s\n", status)

	if errorMsg != "" {
		fmt.Printf("Error: %s\n", errorMsg)
	}

	if resultMap, ok := result["result"].(map[string]interface{}); ok {
		fmt.Println("\n--- Result ---")
		if success, ok := resultMap["success"].(bool); ok {
			fmt.Printf("Success: %v\n", success)
		}
		if commandName, ok := resultMap["command_name"].(string); ok && commandName != "" {
			fmt.Printf("Command: %s\n", commandName)
		}
		if sourceID, ok := resultMap["source_id"].(string); ok && sourceID != "" {
			fmt.Printf("Source ID: %s\n", sourceID)
		}
		if chunks, ok := resultMap["chunks_created"].(float64); ok {
			fmt.Printf("Chunks Created: %.0f\n", chunks)
		}
		if execTime, ok := resultMap["execution_time"].(float64); ok {
			fmt.Printf("Execution Time: %.2fs\n", execTime)
		}
		if execMeta, ok := resultMap["execution_metadata"].(map[string]interface{}); ok {
			if startedAt, ok := execMeta["started_at"].(string); ok && startedAt != "" {
				fmt.Printf("Started At: %s\n", startedAt)
			}
		}
	}

	return nil
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}