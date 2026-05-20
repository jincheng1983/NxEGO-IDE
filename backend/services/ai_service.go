package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AIProvider string

const (
	ProviderLocal  AIProvider = "local"
	ProviderOllama AIProvider = "ollama"
	ProviderCloud  AIProvider = "cloud"
)

type AIConfig struct {
	Provider             AIProvider `json:"provider"`
	LocalModelPath       string     `json:"localModelPath"`
	LocalModelName       string     `json:"localModelName"`
	OllamaURL            string     `json:"ollamaUrl"`
	OllamaModel          string     `json:"ollamaModel"`
	CloudAPIKey          string     `json:"cloudApiKey"`
	CloudAPIURL          string     `json:"cloudApiUrl"`
	CloudModel           string     `json:"cloudModel"`
	EnableGPU            bool       `json:"enableGPU"`
	GPULayers            int        `json:"gpuLayers"`
	MaxTokens            int        `json:"maxTokens"`
	Temperature          float64    `json:"temperature"`
	TopP                 float64    `json:"topP"`
	ContextSize          int        `json:"contextSize"`
	EnableCodeCompletion bool       `json:"enableCodeCompletion"`
	CompletionDelay      int        `json:"completionDelay"`
	EnableChatHistory    bool       `json:"enableChatHistory"`
	SystemPrompt         string     `json:"systemPrompt"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"maxTokens"`
	Temperature float64       `json:"temperature"`
	TopP        float64       `json:"topP"`
	Stream      bool          `json:"stream"`
}

type ChatResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

type CompletionRequest struct {
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Language  string `json:"language"`
	MaxTokens int    `json:"maxTokens"`
}

type CompletionResponse struct {
	Text  string `json:"text"`
	Error string `json:"error,omitempty"`
}

type StreamCallback func(chunk string, done bool, err error)

type AIService struct {
	mu     sync.RWMutex
	config AIConfig
	client *http.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAIService() *AIService {
	ctx, cancel := context.WithCancel(context.Background())
	return &AIService{
		config: AIConfig{
			Provider:     ProviderLocal,
			MaxTokens:    2048,
			Temperature:  0.7,
			TopP:         0.9,
			ContextSize:  32768,
			GPULayers:    28,
			SystemPrompt: "你是易狗 NxEGO IDE 的 AI 编程助手，精通 Go 语言和中文函数映射体系。",
		},
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *AIService) UpdateConfig(config AIConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
}

func (s *AIService) GetConfig() AIConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

func (s *AIService) Chat(messages []ChatMessage) (*ChatResponse, error) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	systemMessages := []ChatMessage{
		{Role: "system", Content: config.SystemPrompt},
	}
	allMessages := append(systemMessages, messages...)

	switch config.Provider {
	case ProviderOllama:
		return s.chatOllama(allMessages, config)
	case ProviderCloud:
		return s.chatCloud(allMessages, config)
	default:
		return s.chatLocal(allMessages, config)
	}
}

func (s *AIService) ChatStream(messages []ChatMessage, callback StreamCallback) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	systemMessages := []ChatMessage{
		{Role: "system", Content: config.SystemPrompt},
	}
	allMessages := append(systemMessages, messages...)

	switch config.Provider {
	case ProviderOllama:
		s.chatOllamaStream(allMessages, config, callback)
	case ProviderCloud:
		s.chatCloudStream(allMessages, config, callback)
	default:
		s.chatLocalStream(allMessages, config, callback)
	}
}

func (s *AIService) Complete(req CompletionRequest) (*CompletionResponse, error) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	if !config.EnableCodeCompletion {
		return &CompletionResponse{Text: ""}, nil
	}

	prompt := fmt.Sprintf(
		"你是一个代码补全助手。请根据以下代码上下文，直接输出补全的代码，不要任何解释。\n\n语言: %s\n\n上下文代码:\n%s",
		req.Language, req.Prefix,
	)

	messages := []ChatMessage{
		{Role: "user", Content: prompt},
	}

	resp, err := s.Chat(messages)
	if err != nil {
		return &CompletionResponse{Error: err.Error()}, err
	}

	return &CompletionResponse{Text: resp.Content}, nil
}

func (s *AIService) Shutdown() {
	s.cancel()
}

func (s *AIService) chatLocal(messages []ChatMessage, config AIConfig) (*ChatResponse, error) {
	return &ChatResponse{
		Content: "本地模型推理功能正在开发中。\n\n当前支持的功能：\n• 代码补全（基于规则）\n• 中英文函数映射查询\n• 项目模板生成\n\n请确保已配置 Qwen2.5-Coder 模型文件路径。",
		Done:    true,
	}, nil
}

func (s *AIService) chatLocalStream(messages []ChatMessage, config AIConfig, callback StreamCallback) {
	placeholder := "本地模型推理功能正在开发中。\n\n当前支持的功能：\n• 代码补全（基于规则）\n• 中英文函数映射查询\n• 项目模板生成\n\n请确保已配置 Qwen2.5-Coder 模型文件路径。"

	for _, r := range placeholder {
		callback(string(r), false, nil)
		time.Sleep(20 * time.Millisecond)
	}
	callback("", true, nil)
}

func (s *AIService) chatOllama(messages []ChatMessage, config AIConfig) (*ChatResponse, error) {
	url := strings.TrimRight(config.OllamaURL, "/") + "/api/chat"

	body := map[string]interface{}{
		"model":    config.OllamaModel,
		"messages": messages,
		"stream":   false,
		"options": map[string]interface{}{
			"temperature": config.Temperature,
			"top_p":       config.TopP,
			"num_predict": config.MaxTokens,
			"num_ctx":     config.ContextSize,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Ollama 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama 返回错误 (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Done bool `json:"done"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &ChatResponse{
		Content: result.Message.Content,
		Done:    result.Done,
	}, nil
}

func (s *AIService) chatOllamaStream(messages []ChatMessage, config AIConfig, callback StreamCallback) {
	url := strings.TrimRight(config.OllamaURL, "/") + "/api/chat"

	body := map[string]interface{}{
		"model":    config.OllamaModel,
		"messages": messages,
		"stream":   true,
		"options": map[string]interface{}{
			"temperature": config.Temperature,
			"top_p":       config.TopP,
			"num_predict": config.MaxTokens,
			"num_ctx":     config.ContextSize,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		callback("", true, fmt.Errorf("序列化请求失败: %w", err))
		return
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		callback("", true, fmt.Errorf("创建请求失败: %w", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		callback("", true, fmt.Errorf("请求 Ollama 失败: %w", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		callback("", true, fmt.Errorf("Ollama 返回错误 (%d): %s", resp.StatusCode, string(respBody)))
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done bool `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if chunk.Done {
			callback("", true, nil)
			return
		}

		if chunk.Message.Content != "" {
			callback(chunk.Message.Content, false, nil)
		}
	}

	if err := scanner.Err(); err != nil {
		callback("", true, fmt.Errorf("读取流式响应失败: %w", err))
		return
	}

	callback("", true, nil)
}

func (s *AIService) chatCloud(messages []ChatMessage, config AIConfig) (*ChatResponse, error) {
	url := strings.TrimRight(config.CloudAPIURL, "/") + "/chat/completions"

	openaiMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	body := map[string]interface{}{
		"model":       config.CloudModel,
		"messages":    openaiMessages,
		"max_tokens":  config.MaxTokens,
		"temperature": config.Temperature,
		"top_p":       config.TopP,
		"stream":      false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.CloudAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求云端 API 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("云端 API 返回错误 (%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("云端 API 返回空响应")
	}

	return &ChatResponse{
		Content: result.Choices[0].Message.Content,
		Done:    true,
	}, nil
}

func (s *AIService) chatCloudStream(messages []ChatMessage, config AIConfig, callback StreamCallback) {
	url := strings.TrimRight(config.CloudAPIURL, "/") + "/chat/completions"

	openaiMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	body := map[string]interface{}{
		"model":       config.CloudModel,
		"messages":    openaiMessages,
		"max_tokens":  config.MaxTokens,
		"temperature": config.Temperature,
		"top_p":       config.TopP,
		"stream":      true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		callback("", true, fmt.Errorf("序列化请求失败: %w", err))
		return
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		callback("", true, fmt.Errorf("创建请求失败: %w", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.CloudAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		callback("", true, fmt.Errorf("请求云端 API 失败: %w", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		callback("", true, fmt.Errorf("云端 API 返回错误 (%d): %s", resp.StatusCode, string(respBody)))
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			callback("", true, nil)
			return
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			if chunk.Choices[0].FinishReason != nil {
				callback("", true, nil)
				return
			}
			if chunk.Choices[0].Delta.Content != "" {
				callback(chunk.Choices[0].Delta.Content, false, nil)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		callback("", true, fmt.Errorf("读取流式响应失败: %w", err))
		return
	}

	callback("", true, nil)
}
