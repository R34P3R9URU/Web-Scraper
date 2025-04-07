package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request returned status %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTML parsing failed: %w", err)
	}

	return doc, nil
}

// Traverse HTML to extract <title>, <h1>-<h3>, and <a href="...">
func traverse(n *html.Node, title *string, headings *[]string, links *[]string) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				*title = n.FirstChild.Data
			}
		case "h1", "h2", "h3":
			text := getTextContent(n)
			if text != "" {
				*headings = append(*headings, fmt.Sprintf("%s: %s", strings.ToUpper(n.Data), text))
			}
		case "a":
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					*links = append(*links, attr.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, title, headings, links)
	}
}

// Helper to extract all text inside an element (e.g., heading)
func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(n)
	return strings.TrimSpace(sb.String())
}

func main() {
	url := "https://runescape.com" // change as needed

	root, err := fetchHTML(url)
	if err != nil {
		log.Fatalf("Error fetching HTML: %v", err)
	}

	var (
		title    string
		headings []string
		links    []string
	)

	traverse(root, &title, &headings, &links)

	fmt.Println("Page Title:", title)
	fmt.Println("\nHeadings:")
	for _, h := range headings {
		fmt.Println(" -", h)
	}

	fmt.Println("\nLinks:")
	for _, link := range links {
		fmt.Println(" -", link)
	}
}
