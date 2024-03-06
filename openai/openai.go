package openai

// Define the message structure.
type ReqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Define the payload structure.
type Payload struct {
	Model          string       `json:"model"`
	Messages       []ReqMessage `json:"messages"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

// Define the struct to match the JSON structure.
type Response struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	Created     int64    `json:"created"`
	Model       string   `json:"model"`
	Choices     []Choice `json:"choices"`
	Usage       Usage    `json:"usage"`
	Fingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
