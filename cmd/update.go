package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CLI to the latest version",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := version

	// Check CLI update
	if err := updateCLI(currentVersion); err != nil {
		return err
	}

	// Check skill version mismatch - just show warning
	if err := checkSkillVersion(currentVersion); err != nil {
		fmt.Fprintf(os.Stderr, "\nNote: %s\n", err.Error())
	}

	return nil
}

func updateCLI(currentVersion string) error {
	latestURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		"darimuri", "open-notebook-cli")

	resp, err := http.Get(latestURL)
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return fmt.Errorf("failed to parse release: %w", err)
	}

	if release.TagName == "" {
		return fmt.Errorf("could not find latest version tag")
	}

	// Remove 'v' prefix if present
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	// Normalize currentVersion by removing 'v' prefix if present
	currentVersionNorm := strings.TrimPrefix(currentVersion, "v")

	if currentVersionNorm != "dev" && currentVersionNorm == latestVersion {
		fmt.Printf("CLI is already on latest version: %s\n", currentVersion)
		return nil
	}

	fmt.Printf("Updating CLI from %s to %s...\n", currentVersion, release.TagName)

	// Build download URL
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}
	filename := fmt.Sprintf("open-notebook-%s-%s%s", goos, goarch, ext)
	downloadURL := fmt.Sprintf("https://github.com/darimuri/open-notebook-cli/releases/download/%s/%s",
		release.TagName, filename)

	// Download the binary
	fmt.Printf("Downloading %s...\n", downloadURL)

	httpResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return fmt.Errorf("failed to download binary: %s", httpResp.Status)
	}

	// Determine install path
	installPath := "/usr/local/bin/open-notebook"
	if goos == "windows" {
		installPath = "open-notebook.exe"
	}

	// Check if we can write to /usr/local/bin
	if _, err := os.Stat("/usr/local/bin"); err == nil {
		if os.WriteFile(installPath, nil, 0755) != nil {
			// Try user bin
			homeDir, _ := os.UserHomeDir()
			installPath = homeDir + "/.local/bin/open-notebook" + ext
		}
	} else {
		homeDir, _ := os.UserHomeDir()
		installPath = homeDir + "/.local/bin/open-notebook" + ext
	}

	// Save to temp file first
	tmpPath := installPath + ".tmp"
	outFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(outFile, httpResp.Body)
	outFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	os.Chmod(tmpPath, 0755)
	os.Rename(tmpPath, installPath)

	fmt.Printf("CLI updated to %s. Installed to: %s\n", latestVersion, installPath)
	return nil
}

func checkSkillVersion(cliVersion string) error {
	skillURL := "https://raw.githubusercontent.com/darimuri/open-notebook-cli/main/skills/open-notebook/SKILL.md"

	resp, err := http.Get(skillURL)
	if err != nil {
		return fmt.Errorf("failed to fetch skill: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get skill: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read skill: %w", err)
	}

	// Parse version from SKILL.md (look for version: "x.y.z" in frontmatter)
	versionRe := regexp.MustCompile(`version:\s*"?([0-9]+\.[0-9]+\.[0-9]+)"?`)
	matches := versionRe.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		// No version in skill, skip check
		return nil
	}

	skillVersion := matches[1]

	// Normalize versions by removing 'v' prefix
	cliVersionNorm := strings.TrimPrefix(cliVersion, "v")
	skillVersionNorm := strings.TrimPrefix(skillVersion, "v")

	// Compare versions
	if cliVersionNorm != "dev" && cliVersionNorm != skillVersionNorm {
		return fmt.Errorf("skill version (%s) does not match CLI version (%s). Plugin update required: https://github.com/darimuri/open-notebook-cli", skillVersion, cliVersion)
	}

	return nil
}