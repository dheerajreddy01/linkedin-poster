package news

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"linkedin-poster/internal/models"
)

type Fetcher struct{}

func New() *Fetcher { return &Fetcher{} }

type RSS struct {
	Channel struct {
		Items []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
			Desc    string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type NewsAPIResponse struct {
	Articles []struct {
		Title       string    `json:"title"`
		URL         string    `json:"url"`
		Source      struct{ Name string `json:"name"` } `json:"source"`
		PublishedAt time.Time `json:"publishedAt"`
		Description string    `json:"description"`
	} `json:"articles"`
}

// Tech News & AI sources only
var sources = []struct {
	Name  string
	URL   string
	Topic string
}{
	// AI & LLMs
	{"OpenAI Blog",       "https://openai.com/blog/rss/",                             "AI & LLMs"},
	{"Google AI Blog",    "https://blog.google/technology/ai/rss/",                   "AI & LLMs"},
	{"Hugging Face",      "https://huggingface.co/blog/feed.xml",                     "AI & LLMs"},
	{"MIT AI News",       "https://news.mit.edu/topic/artificial-intelligence2/feed", "AI & LLMs"},
	{"r/MachineLearning", "https://www.reddit.com/r/MachineLearning/.rss?limit=5",    "AI & LLMs"},
	{"r/artificial",      "https://www.reddit.com/r/artificial/.rss?limit=5",         "AI & LLMs"},
	// Tech News
	{"Hacker News",       "https://news.ycombinator.com/rss",                         "Tech News"},
	{"TechCrunch",        "https://techcrunch.com/feed/",                             "Tech News"},
	{"Ars Technica",      "https://feeds.arstechnica.com/arstechnica/technology-lab", "Tech News"},
	{"The Verge",         "https://www.theverge.com/tech/rss/index.xml",              "Tech News"},
	{"Wired",             "https://www.wired.com/feed/rss",                           "Tech News"},
	{"dev.to",            "https://dev.to/feed",                                      "Tech News"},
	{"InfoQ",             "https://feed.infoq.com",                                   "Tech News"},
	{"r/technology",      "https://www.reddit.com/r/technology/.rss?limit=5",         "Tech News"},
}

var newsAPIQueries = []struct{ Query, Topic string }{
	{"artificial intelligence LLM GPT Claude Gemini 2025", "AI & LLMs"},
	{"machine learning deep learning neural network 2025",  "AI & LLMs"},
	{"tech news software engineering cloud AWS 2025",       "Tech News"},
}

func (f *Fetcher) FetchAll(apiKey string) []models.NewsItem {
	var all []models.NewsItem
	seen := map[string]bool{}

	for _, src := range sources {
		items := f.fetchRSS(src.URL, src.Name, src.Topic)
		for _, item := range items {
			if !seen[item.URL] && item.URL != "" {
				seen[item.URL] = true
				all = append(all, item)
			}
		}
	}

	if apiKey != "" {
		for _, q := range newsAPIQueries {
			items := f.fetchNewsAPI(apiKey, q.Query, q.Topic)
			for _, item := range items {
				if !seen[item.URL] && item.URL != "" {
					seen[item.URL] = true
					all = append(all, item)
				}
			}
		}
	}

	all = filterQuality(all)
	log.Printf("✅ Fetched %d tech/AI news items", len(all))
	return all
}

func (f *Fetcher) fetchRSS(feedURL, sourceName, topic string) []models.NewsItem {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", feedURL, nil)
	req.Header.Set("User-Agent", "PostPilot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("RSS error [%s]: %v", sourceName, err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil
	}

	var items []models.NewsItem
	for i, item := range rss.Channel.Items {
		if i >= 4 { break }
		if item.Title == "" || item.Link == "" { continue }

		pubAt, _ := time.Parse(time.RFC1123Z, item.PubDate)
		if pubAt.IsZero() {
			pubAt, _ = time.Parse(time.RFC1123, item.PubDate)
		}
		if pubAt.IsZero() {
			pubAt = time.Now()
		}

		summary := stripHTML(item.Desc)
		if len(summary) > 300 {
			summary = summary[:300] + "..."
		}

		items = append(items, models.NewsItem{
			Title:       cleanTitle(item.Title),
			URL:         item.Link,
			Source:      sourceName,
			Topic:       topic,
			Summary:     summary,
			PublishedAt: pubAt,
		})
	}
	return items
}

func (f *Fetcher) fetchNewsAPI(apiKey, query, topic string) []models.NewsItem {
	endpoint := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&sortBy=publishedAt&pageSize=5&language=en&apiKey=%s",
		url.QueryEscape(query), apiKey,
	)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result NewsAPIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	var items []models.NewsItem
	for _, a := range result.Articles {
		if a.Title == "" || a.Title == "[Removed]" { continue }
		items = append(items, models.NewsItem{
			Title:       cleanTitle(a.Title),
			URL:         a.URL,
			Source:      a.Source.Name,
			Topic:       topic,
			Summary:     a.Description,
			PublishedAt: a.PublishedAt,
		})
	}
	return items
}

func filterQuality(items []models.NewsItem) []models.NewsItem {
	skip := []string{"sponsored", "advertisement", "buy now", "giveaway", "subscribe now"}
	var out []models.NewsItem
	for _, item := range items {
		title := strings.ToLower(item.Title)
		bad := false
		for _, w := range skip {
			if strings.Contains(title, w) { bad = true; break }
		}
		if !bad && len(item.Title) > 10 {
			out = append(out, item)
		}
	}
	return out
}

func stripHTML(s string) string {
	inTag := false
	var b strings.Builder
	for _, ch := range s {
		if ch == '<' { inTag = true; continue }
		if ch == '>' { inTag = false; continue }
		if !inTag { b.WriteRune(ch) }
	}
	return strings.TrimSpace(b.String())
}

func cleanTitle(s string) string {
	s = stripHTML(s)
	s = strings.TrimSpace(s)
	if idx := strings.LastIndex(s, " | "); idx > 20 { s = s[:idx] }
	if idx := strings.LastIndex(s, " - "); idx > 20 { s = s[:idx] }
	return s
}
