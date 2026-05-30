package crawler

import (
	"html"
	"net/url"
	"strings"

	nethtml "golang.org/x/net/html"
)

type Link struct {
	URL        string
	Text       string
	IsInternal bool
}

func ExtractLinks(baseURL string, htmlContent string) ([]Link, error) {
	doc, err := nethtml.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	var links []Link
	extractLinks(doc, base, &links)

	return links, nil
}

var ignoredExtensions = []string{
	".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico",
	".pdf", ".zip", ".tar", ".gz", ".rar",
	".css", ".js", ".map",
	".mp4", ".avi", ".mov", ".mp3", ".wav",
}

func isAssetURL(rawURL string) bool {
	rawURL = strings.ToLower(rawURL)
	for _, ext := range ignoredExtensions {
		if strings.HasSuffix(rawURL, ext) {
			return true
		}
	}
	return false
}

func extractLinks(n *nethtml.Node, base *url.URL, links *[]Link) {
	if n.Type == nethtml.ElementNode && (n.Data == "a" || n.Data == "link") {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				href := strings.TrimSpace(attr.Val)
				if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") {
					continue
				}

				parsedURL, err := url.Parse(href)
				if err != nil {
					continue
				}

				// Resolve relative URLs
				absoluteURL := base.ResolveReference(parsedURL).String()

				// Normalize URL (remove fragments, trailing slashes)
				absoluteURL = normalizeURL(absoluteURL)

				// Skip asset URLs (images, PDFs, etc.)
				if isAssetURL(absoluteURL) {
					continue
				}

				// Check if internal (same domain)
				isInternal := parsedURL.Host == "" || parsedURL.Host == base.Host

				// Skip non-http(s) schemes (mailto, tel, etc.)
				scheme := strings.ToLower(parsedURL.Scheme)
				if scheme != "" && scheme != "http" && scheme != "https" {
					continue
				}

				*links = append(*links, Link{
					URL:        absoluteURL,
					Text:       getLinkText(n),
					IsInternal: isInternal,
				})
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, base, links)
	}
}

func getLinkText(n *nethtml.Node) string {
	if n.Type == nethtml.TextNode {
		return strings.TrimSpace(html.UnescapeString(n.Data))
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getLinkText(c)
	}
	return strings.TrimSpace(text)
}

func normalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	// Remove fragment
	u.Fragment = ""

	// Remove trailing slash if path has multiple segments
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// Lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	if u.Host != "" {
		u.Host = strings.ToLower(u.Host)
	}

	return u.String()
}

func NormalizeURL(rawURL string) string {
	return normalizeURL(rawURL)
}