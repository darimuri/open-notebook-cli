package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage and invoke skills",
	Long:  `List and invoke Claude Code skills`,
}

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available skills",
	RunE:  runSkillsList,
}

var skillsInvokeCmd = &cobra.Command{
	Use:   "invoke [skill_name]",
	Short: "Invoke a skill",
	Args:  cobra.ExactArgs(1),
	RunE:  runSkillsInvoke,
}

func init() {
	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsInvokeCmd)
	rootCmd.AddCommand(skillsCmd)
}

func runSkillsList(cmd *cobra.Command, args []string) error {
	fmt.Println("Available skills:")
	fmt.Println("  brainstorming - Use brainstorming workflow")
	fmt.Println("  debugging    - Use debugging workflow")
	return nil
}

func runSkillsInvoke(cmd *cobra.Command, args []string) error {
	skillName := args[0]
	fmt.Printf("Invoking skill: %s\n", skillName)
	fmt.Println("Note: Claude Code skills integration requires running within Claude Code")
	os.Exit(0)
	return nil
}