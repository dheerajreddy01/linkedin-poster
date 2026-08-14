package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"linkedin-poster/internal/models"
)

// Groq exposes an OpenAI-compatible Chat Completions API at this base URL.
const groqBaseURL = "https://api.groq.com/openai/v1"

// defaultModel is a fast, high-quality Groq-hosted model. Override with GROQ_MODEL.
const defaultModel = "llama-3.3-70b-versatile"

type Service struct {
	client *openai.Client
	model  string
}

func New() *Service {
	config := openai.DefaultConfig(os.Getenv("GROQ_API_KEY"))
	config.BaseURL = groqBaseURL

	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = defaultModel
	}

	return &Service{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}
}

// GeneratePost creates a LinkedIn post from a news item
func (s *Service) GeneratePost(item models.NewsItem, authorName string) (string, error) {
	tone := toneForTopic(item.Topic)

	prompt := fmt.Sprintf(`You are %s, a Senior Software Engineer at Capital One specializing in Go, Java, Python, AWS, and Data Science.

Write a LinkedIn post about this tech news. Make it feel authentic, personal, and insightful — not corporate or robotic.

News Title: %s
Source: %s
Topic: %s
Summary: %s

Tone: %s

Rules:
- Start with a strong hook (question, bold statement, or surprising fact)
- 150-250 words max
- Share YOUR perspective and what this means for engineers/developers
- Add 1-2 practical takeaways
- End with a thought-provoking question to drive comments
- Add 3-5 relevant hashtags at the end
- Do NOT use "Excited to share" or "I'm thrilled" — be genuine
- Do NOT use excessive emojis — max 2-3 total
- Reference your experience with %s where relevant

Write only the post content, nothing else.`,
		authorName,
		item.Title,
		item.Source,
		item.Topic,
		item.Summary,
		tone,
		item.Topic,
	)

	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       s.model,
		Temperature: 0.8,
		MaxTokens:   500,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("groq error: %v", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// RegeneratePost rewrites a post with a different angle
func (s *Service) RegeneratePost(original string, instruction string) (string, error) {
	prompt := fmt.Sprintf(`Rewrite this LinkedIn post with the following instruction: %s

Original post:
%s

Keep the same news topic but apply the instruction. Same length (150-250 words). Same hashtags or update them if needed.
Write only the post content, nothing else.`,
		instruction, original,
	)

	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       s.model,
		Temperature: 0.9,
		MaxTokens:   500,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("groq error: %v", err)
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func toneForTopic(topic string) string {
	switch topic {
	case "Go / Backend Engineering":
		return "technical but approachable, like a senior engineer sharing hard-won lessons"
	case "Data Science & ML":
		return "curious and analytical, like someone actively learning and experimenting"
	case "AWS & Cloud":
		return "practical and cost-conscious, like someone who's debugged production AWS issues at 2am"
	case "AI & LLMs":
		return "thoughtful and forward-looking, neither hype nor dismissive"
	case "Open Source":
		return "community-oriented and enthusiastic about collaboration"
	default:
		return "conversational, honest, and relatable to software engineers"
	}
}
