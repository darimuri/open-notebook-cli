package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to the latest version",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := version

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

	// Simple JSON parsing - extract tag_name
	lines := strings.Split(string(body), "\n")
	var latestTag string
	for _, line := range lines {
		if strings.Contains(line, `"tag_name"`) {
			parts := strings.Split(line, `"`)
			if len(parts) >= 4 {
				latestTag = parts[3]
			}
			break
		}
	}

	if latestTag == "" {
		return fmt.Errorf("could not find latest version tag")
	}

	// Remove 'v' prefix if present
	latestVersion := strings.TrimPrefix(latestTag, "v")

	if currentVersion != "dev" && currentVersion == latestVersion {
		fmt.Printf("Already on latest version: %s\n", currentVersion)
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", currentVersion, latestVersion)

	// Build download URL
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}
	filename := fmt.Sprintf("open-notebook-%s-%s%s", goos, goarch, ext)
	downloadURL := fmt.Sprintf("https://github.com/darimuri/open-notebook-cli/releases/download/%s/%s",
		latestTag, filename)

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

	fmt.Printf("Updated to %s. Installed to: %s\n", latestVersion, installPath)
	return nil
}