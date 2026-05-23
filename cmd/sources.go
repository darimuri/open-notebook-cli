package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/darimuri/open-notebook-cli/internal/api"
	"github.com/darimuri/open-notebook-cli/internal/crawler"
)

var (
	recursive      bool
	maxDepth       int
	maxCount       int
	pollingPeriod  int
	sourceNotebook string
	sourceFile     string
	skipEmbed      bool
	pageFlag       int
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Manage sources",
	Long:  `List, add, upload, download, and manage sources`,
}

var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sources",
	RunE:  runSourcesList,
}

var sourcesAddCmd = &cobra.Command{
	Use:   "add [url/file/text...]",
	Short: "Add source(s) - URL, file, or text",
	Long: `Add source(s) to notebooks.

Examples:
  # Add a single URL
  open-notebook sources add https://example.com/article

  # Add multiple URLs
  open-notebook sources add https://site.com/page1 https://site.com/page2

  # Add with recursive crawling (all internal links)
  open-notebook sources add -r https://site.com/docs

  # Add with depth limit
  open-notebook sources add -r --depth 3 https://site.com/docs

  # Add from file (one URL per line)
  open-notebook sources add --file urls.txt

  # Add text content
  open-notebook sources add --text "Some important notes"

  # Add to specific notebook
  open-notebook sources add -n notebook-id https://example.com`,
	RunE: runSourcesAdd,
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

var sourcesInsightsCmd = &cobra.Command{
	Use:   "insights [source_id]",
	Short: "Get source insights",
	Args:  cobra.ExactArgs(1),
	RunE:  runSourcesInsights,
}

var sourcesStatusCmd = &cobra.Command{
	Use:   "status [source_id]",
	Short: "Get source processing status",
	Args:  cobra.ExactArgs(1),
	RunE:  runSourcesStatus,
}

var sourcesEmbedCmd = &cobra.Command{
	Use:   "embed [source_id]",
	Short: "Embed a source for vector search",
	Long: `Embed a source for vector search.

Use --wait to monitor the embed job until completion.`,
	Args: cobra.ExactArgs(1),
	RunE:  runSourcesEmbed,
}

var embedWait bool

var sourcesEmbedBatchCmd = &cobra.Command{
	Use:   "embed-batch",
	Short: "Embed all non-embedded sources",
	Long: `Embed all non-embedded sources in all pages.

This command will:
	1. List all sources across all pages
	2. Filter for non-embedded sources (embedded_chunks=0, status=completed)
	3. Trigger embed for each source
	4. Monitor until all complete
	5. Stop and report if any embed fails`,
	RunE: runSourcesEmbedBatch,
}

func init() {
	sourcesCmd.AddCommand(sourcesListCmd)
	sourcesCmd.AddCommand(sourcesAddCmd)
	sourcesCmd.AddCommand(sourcesUploadCmd)
	sourcesCmd.AddCommand(sourcesDownloadCmd)
	sourcesCmd.AddCommand(sourcesRetryCmd)
	sourcesCmd.AddCommand(sourcesInsightsCmd)
	sourcesCmd.AddCommand(sourcesStatusCmd)
	sourcesCmd.AddCommand(sourcesEmbedCmd)
	sourcesCmd.AddCommand(sourcesEmbedBatchCmd)
	rootCmd.AddCommand(sourcesCmd)

	// Add command flags
	sourcesListCmd.Flags().StringVarP(&sourceNotebook, "notebook", "n", "", "Filter by notebook ID")
	sourcesListCmd.Flags().IntVar(&pageFlag, "page", 1, "Page number (default 1)")

	sourcesAddCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively crawl internal links")
	sourcesAddCmd.Flags().IntVar(&maxDepth, "depth", 0, "Maximum crawl depth (0 = unlimited)")
	sourcesAddCmd.Flags().StringVarP(&sourceNotebook, "notebook", "n", "", "Notebook ID to add sources to")
	sourcesAddCmd.Flags().StringVarP(&sourceFile, "file", "f", "", "Read URLs from file (one per line)")
	sourcesAddCmd.Flags().String("text", "", "Add text content as source")
	sourcesAddCmd.Flags().BoolVar(&skipEmbed, "skip-embed", false, "Skip embedding (default: embed)")

	sourcesEmbedBatchCmd.Flags().StringVarP(&sourceNotebook, "notebook", "n", "", "Notebook ID to filter sources")
	sourcesEmbedBatchCmd.Flags().IntVar(&maxCount, "max", 0, "Maximum number of sources to embed (0 = all)")
	sourcesEmbedBatchCmd.Flags().IntVar(&pollingPeriod, "polling-period", 10, "Polling interval in seconds")

	sourcesEmbedCmd.Flags().BoolVar(&embedWait, "wait", false, "Wait for embed to complete")
	sourcesEmbedCmd.Flags().IntVar(&pollingPeriod, "polling-period", 10, "Polling interval in seconds")
}

func runSourcesList(cmd *cobra.Command, args []string) error {
	client := getClient()

	url := "/api/sources"
	page := pageFlag
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * 50
	if sourceNotebook != "" {
		url = fmt.Sprintf("/api/sources?notebook=%s&offset=%d", sourceNotebook, offset)
	} else {
		url = fmt.Sprintf("/api/sources?offset=%d", offset)
	}

	var sources []api.SourceResponse
	err := client.Get(url, &sources)
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	return outputJSON(sources)
}

func runSourcesAdd(cmd *cobra.Command, args []string) error {
	client := getClient()

	// Check for --text flag
	if text := cmd.Flag("text").Value.String(); text != "" {
		return addTextSource(client, text)
	}

	// Check for --file flag
	filePath := sourceFile
	if filePath != "" {
		return addSourcesFromFile(client, filePath)
	}

	if len(args) == 0 {
		return fmt.Errorf("no sources specified. Use add [urls], --text, or --file")
	}

	// Collect URLs to add
	var urls []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			urls = append(urls, arg)
		} else if _, err := os.Stat(arg); err == nil {
			// It's a file - read URLs from it
			fileURLs, err := readURLsFromFile(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to read file %s: %v\n", arg, err)
			} else {
				urls = append(urls, fileURLs...)
			}
		} else {
			return fmt.Errorf("invalid source: %s (not a URL or file)", arg)
		}
	}

	if len(urls) == 0 {
		return fmt.Errorf("no valid URLs to add")
	}

	// If recursive, crawl each URL
	if recursive {
		return addSourcesRecursive(client, urls)
	}

	// Add URLs directly
	return addSources(client, urls)
}

func addTextSource(client *api.Client, text string) error {
	req := api.SourceCreate{
		Type:    "text",
		Content: text,
	}

	if sourceNotebook != "" {
		req.Notebooks = []string{sourceNotebook}
	}

	embed := !skipEmbed
	req.Embed = &embed

	var result api.SourceResponse
	err := client.Post("/api/sources/json", req, &result)
	if err != nil {
		return fmt.Errorf("failed to add text source: %w", err)
	}

	fmt.Printf("Added text source: %s\n", result.ID)
	return nil
}

func addSourcesFromFile(client *api.Client, filePath string) error {
	urls, err := readURLsFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read URLs from file: %w", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found in file")
	}

	if recursive {
		return addSourcesRecursive(client, urls)
	}

	return addSources(client, urls)
}

func readURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
			urls = append(urls, url)
		}
	}

	return urls, scanner.Err()
}

func addSources(client *api.Client, urls []string) error {
	added := 0
	failed := 0

	for _, url := range urls {
		req := api.SourceCreate{
			Type: "link",
			URL:  url,
		}

		if sourceNotebook != "" {
			req.Notebooks = []string{sourceNotebook}
		}

		embed := !skipEmbed
		req.Embed = &embed

		var result api.SourceResponse
		err := client.Post("/api/sources/json", req, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add %s: %v\n", url, err)
			failed++
		} else {
			fmt.Printf("Added: %s\n", url)
			added++
		}
	}

	fmt.Printf("\nAdded: %d, Failed: %d\n", added, failed)
	return nil
}

func addSourcesRecursive(client *api.Client, startURLs []string) error {
	visited := make(map[string]bool)
	var queue []string

	// Initialize queue with starting URLs
	for _, url := range startURLs {
		normalized := crawler.NormalizeURL(url)
		if !visited[normalized] {
			visited[normalized] = true
			queue = append(queue, normalized)
		}
	}

	depth := 0
	added := 0
	failed := 0

	for len(queue) > 0 {
		if maxDepth > 0 && depth >= maxDepth {
			fmt.Printf("Reached max depth %d, stopping\n", maxDepth)
			break
		}

		url := queue[0]
		queue = queue[1:]

		fmt.Printf("[%d] Crawling: %s\n", depth, url)

		// Fetch page
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to fetch: %v\n", err)
			failed++
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to read body: %v\n", err)
			failed++
			continue
		}

		// Add as source
		req := api.SourceCreate{
			Type: "link",
			URL:  url,
		}
		if sourceNotebook != "" {
			req.Notebooks = []string{sourceNotebook}
		}
		embed := !skipEmbed
		req.Embed = &embed

		var result api.SourceResponse
		err = client.Post("/api/sources/json", req, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to add source: %v\n", err)
			failed++
		} else {
			fmt.Printf("  Added: %s (id: %s)\n", url, result.ID)
			added++
		}

		// Extract internal links
		links, err := crawler.ExtractLinks(url, string(body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to extract links: %v\n", err)
			continue
		}

		// Add internal links to queue
		for _, link := range links {
			if link.IsInternal && !visited[link.URL] {
				visited[link.URL] = true
				queue = append(queue, link.URL)
			}
		}

		// Progress update
		fmt.Printf("  Queue: %d, Visited: %d\n", len(queue), len(visited))
	}

	fmt.Printf("\nCrawl complete! Added: %d, Failed: %d, Visited: %d\n", added, failed, len(visited))
	return nil
}

func runSourcesUpload(cmd *cobra.Command, args []string) error {
	client := getClient()

	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	req := map[string]any{
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

func runSourcesInsights(cmd *cobra.Command, args []string) error {
	client := getClient()

	var result any
	err := client.Get("/api/sources/"+args[0]+"/insights", &result)
	if err != nil {
		return fmt.Errorf("failed to get insights: %w", err)
	}

	return outputJSON(result)
}

func runSourcesStatus(cmd *cobra.Command, args []string) error {
	client := getClient()

	var result api.SourceStatusResponse
	err := client.Get("/api/sources/"+args[0]+"/status", &result)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	return outputJSON(result)
}

func runSourcesEmbed(cmd *cobra.Command, args []string) error {
	client := getClient()

	req := api.EmbedRequest{
		ItemID:   args[0],
		ItemType: "source",
	}

	var result api.EmbedResponse
	err := client.Post("/api/embed", req, &result)
	if err != nil {
		return fmt.Errorf("failed to embed source: %w", err)
	}

	fmt.Printf("[%s] - Started embed (cmd: %s)\n", result.ItemID, result.CommandID)

	if !embedWait {
		fmt.Printf("Embedded: %s (%s)\n", result.ItemID, result.Message)
		return nil
	}

	// Wait for embed to complete
	fmt.Println("\nMonitoring embed progress...")
	pollingInterval := 10
	if pollingPeriod > 0 {
		pollingInterval = pollingPeriod
	}
	startTime := time.Now()
	sourceID := result.ItemID
	cmdID := result.CommandID

	for {
		elapsed := time.Since(startTime).Round(time.Second)

		var cmdResult api.CommandJobStatus
		err := client.Get("/api/commands/jobs/"+cmdID, &cmdResult)
		if err != nil {
			return fmt.Errorf("failed to get command status: %w", err)
		}

		switch cmdResult.Status {
		case "completed":
			chunks := 0
			if cmdResult.Result != nil {
				if m, ok := cmdResult.Result.(map[string]any); ok {
					if c, ok := m["chunks_created"].(float64); ok {
						chunks = int(c)
					}
				}
			}
			fmt.Printf("[%s] (chunks: %d, DONE, elapsed: %s)\n", sourceID, chunks, elapsed)
			fmt.Printf("\nEmbed completed successfully! (%s)\n", elapsed)
			return nil
		case "failed":
			fmt.Printf("[%s] - FAILED: %s\n", sourceID, cmdResult.ErrorMessage)
			return fmt.Errorf("embed failed for source %s: %s", sourceID, cmdResult.ErrorMessage)
		case "running", "pending", "new", "":
			fmt.Printf("[%s] (cmd: %s, status: %s, elapsed: %s)\n", sourceID, cmdID, cmdResult.Status, elapsed)
		}

		time.Sleep(time.Duration(pollingInterval) * time.Second)
	}
}

func runSourcesEmbedBatch(cmd *cobra.Command, args []string) error {
	client := getClient()

	// Step 1: Get non-embedded sources (paginated, stop early if max reached)
	fmt.Println("Fetching sources...")
	var allSources []api.SourceResponse
	var nonEmbedded []api.SourceResponse
	page := 1
	for {
		var sources []api.SourceResponse
		var err error
		offset := (page - 1) * 50
		if sourceNotebook != "" {
			err = client.Get(fmt.Sprintf("/api/sources?notebook=%s&offset=%d", sourceNotebook, offset), &sources)
		} else {
			err = client.Get(fmt.Sprintf("/api/sources?offset=%d", offset), &sources)
		}
		if err != nil {
			return fmt.Errorf("failed to fetch sources page %d: %w", page, err)
		}
		if len(sources) == 0 {
			break
		}

		// Filter non-embedded from this page
		for _, s := range sources {
			if s.EmbeddedChunks == 0 && s.Status == "completed" {
				nonEmbedded = append(nonEmbedded, s)
			}
		}
		allSources = append(allSources, sources...)
		fmt.Printf("Page %d: %d sources, %d non-embedded (total: %d, non-embedded: %d)\n",
			page, len(sources), len(nonEmbedded), len(allSources), len(nonEmbedded))

		// Early exit if we have enough non-embedded sources
		if maxCount > 0 && len(nonEmbedded) >= maxCount {
			fmt.Printf("Reached max limit %d, stopping page fetch\n", maxCount)
			break
		}
		page++
	}
	fmt.Printf("\nTotal sources: %d, Non-embedded: %d\n", len(allSources), len(nonEmbedded))

	if len(nonEmbedded) == 0 {
		fmt.Println("All sources are already embedded!")
		return nil
	}

	// Apply max limit if specified
	if maxCount > 0 && len(nonEmbedded) > maxCount {
		nonEmbedded = nonEmbedded[:maxCount]
		fmt.Printf("Limited to %d sources\n", maxCount)
	}

	// Step 2: Sequential embed - one at a time
	fmt.Println("\nStarting sequential embed...")
	for i, s := range nonEmbedded {
		fmt.Printf("[%d/%d] %s - Embedding...\n", i+1, len(nonEmbedded), s.Title)

		// Trigger embed
		req := api.EmbedRequest{
			ItemID:   s.ID,
			ItemType: "source",
		}
		var result api.EmbedResponse
		err := client.Post("/api/embed", req, &result)
		if err != nil {
			return fmt.Errorf("failed to trigger embed for %s: %w", s.ID, err)
		}

		// Poll for completion
		pollingInterval := 10
		if pollingPeriod > 0 {
			pollingInterval = pollingPeriod
		}
		startTime := time.Now()
		for {
			var cmdResult api.CommandJobStatus
			err := client.Get("/api/commands/jobs/"+result.CommandID, &cmdResult)
			if err != nil {
				return fmt.Errorf("failed to get command status for %s: %w", s.ID, err)
			}

			elapsed := time.Since(startTime).Round(time.Second)

			switch cmdResult.Status {
			case "completed":
				chunks := 0
				if cmdResult.Result != nil {
					if m, ok := cmdResult.Result.(map[string]any); ok {
						if c, ok := m["chunks_created"].(float64); ok {
							chunks = int(c)
						}
					}
				}
				fmt.Printf("[%s] DONE (chunks: %d, elapsed: %s)\n", s.Title, chunks, elapsed)
			case "failed":
				return fmt.Errorf("embed failed for source %s: %s", s.ID, cmdResult.ErrorMessage)
			case "running", "pending", "new", "":
				fmt.Printf("[%s] status: %s, elapsed: %s (polling every %ds)\n", s.Title, cmdResult.Status, elapsed, pollingInterval)
				time.Sleep(time.Duration(pollingInterval) * time.Second)
				continue
			}

			break
		}
	}

	fmt.Printf("\nAll %d sources embedded successfully!\n", len(nonEmbedded))
	return nil
}