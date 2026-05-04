// Package agents defines the base Agent interface and shared LLM client
// used by all specialist agents in the pipeline.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"log"
)

// ============================================================
// AGENT INTERFACE — every specialist agent implements this
// ============================================================

type Agent interface {
	Name() string
	Run(ctx context.Context, arch *EnrichedArchitecture) (*EnrichedArchitecture, error)
}

// ============================================================
// LLM CLIENT — standard HTTP wrapper for OpenRouter
// ============================================================

type LLMClient struct {
	model  string
	apiKey string
	client *http.Client
}

func NewLLMClient() *LLMClient {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "openai/gpt-oss-120b:free"
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is not set. Please set it before running.")
	}

	return &LLMClient{
		model:  model,
		apiKey: apiKey,
		client: &http.Client{Timeout: 300 * time.Second},
	}
}

// Invoke sends a prompt to OpenRouter using standard HTTP requests
// systemPrompt   = the agent's persona and instructions
// userPrompt     = the actual task with context injected
func (c *LLMClient) Invoke(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	var messages []Message
	if systemPrompt != "" {
		messages = append(messages, Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, Message{Role: "user", Content: userPrompt})

	reqBody := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
		"reasoning": map[string]bool{
			"enabled": true,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	var resp *http.Response
	var respBytes []byte
	maxRetries := 5

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(bodyBytes))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("HTTP-Referer", "https://github.com/mo7amedgom3a/arch-visualizer") // Recommended for OpenRouter
		req.Header.Set("X-Title", "Arch-Visualizer")                                      // Recommended for OpenRouter

		resp, err = c.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("http request failed: %w", err)
		}

		if resp.StatusCode == 429 {
			bodyText, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if attempt < maxRetries {
				delay := time.Duration(10*(attempt+1)) * time.Second
				log.Printf("Received 429 Rate Limit from OpenAI: %s\nRetrying in %v (Attempt %d/%d)...", string(bodyText), delay, attempt+1, maxRetries)
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("API returned 429 too many times: %s", string(bodyText))
		}

		respBytes, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("API returned unexpected status code: %d: %s", resp.StatusCode, string(respBytes))
		}
		
		if len(respBytes) == 0 {
			return "", fmt.Errorf("API returned empty response body with status 200")
		}
		
		break // Success, exit retry loop
	}

	var respData struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				Reasoning string `json:"reasoning,omitempty"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBytes, &respData); err != nil {
		log.Printf("[LLMClient] Raw response that failed to decode: %s", string(respBytes))
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(respData.Choices) == 0 {
		return "", fmt.Errorf("empty response from LLM: %s", string(respBytes))
	}

	return respData.Choices[0].Message.Content, nil
}

// ============================================================
// HELPERS — shared across agents
// ============================================================

// ArchToJSON serializes the architecture for injection into prompts.
func ArchToJSON(arch *Architecture) string {
	b, _ := json.MarshalIndent(arch, "", "  ")
	return string(b)
}

// ParseJSONBlock extracts a JSON block from LLM output that may contain
// surrounding markdown fences or explanation text.
func ParseJSONBlock(raw string) ([]byte, error) {
	// Find first '{' or '['
	start := -1
	for i, ch := range raw {
		if ch == '{' || ch == '[' {
			start = i
			break
		}
	}
	if start == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	// Find last '}' or ']'
	end := -1
	for i := len(raw) - 1; i >= 0; i-- {
		if raw[i] == '}' || raw[i] == ']' {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, fmt.Errorf("no closing bracket found in response")
	}

	return []byte(raw[start : end+1]), nil
}
