package tokens

import (
	"fmt"
	"io"

	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

// Parser extracts token usage from provider responses
type Parser struct {
	parsers map[string]ProviderParser
}

// NewParser creates a new response parser with registered provider parsers
func NewParser() *Parser {
	return &Parser{
		parsers: map[string]ProviderParser{
			"openai":      &OpenAIParser{},
			"azureopenai": &OpenAIParser{}, // Azure uses same format as OpenAI
			"anthropic":   &AnthropicParser{},
			"cohere":      &CohereParser{},
			"google":      &GoogleParser{},
		},
	}
}

// RegisterParser registers a custom provider parser
func (p *Parser) RegisterParser(provider string, parser ProviderParser) {
	p.parsers[provider] = parser
}

// ParseResponse extracts token usage from a response body
func (p *Parser) ParseResponse(provider string, body []byte) (*model.TokenUsage, error) {
	parser, ok := p.parsers[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return parser.ParseResponse(body)
}

// ParseResponseReader is a convenience method that reads from an io.Reader
func (p *Parser) ParseResponseReader(provider string, reader io.Reader) (*model.TokenUsage, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return p.ParseResponse(provider, body)
}

// ParseStreamedResponse handles streaming responses where usage may be in the final chunk
func (p *Parser) ParseStreamedResponse(provider string, chunks [][]byte) (*model.TokenUsage, error) {
	// For streaming, usage is typically in the last chunk
	// Try parsing chunks in reverse order

	for i := len(chunks) - 1; i >= 0; i-- {
		chunk := chunks[i]

		// Try to parse this chunk
		usage, err := p.ParseResponse(provider, chunk)
		if err == nil {
			return usage, nil
		}
	}

	return nil, fmt.Errorf("no usage data found in streamed response")
}
