package api

import (
	"time"
)

type Trace struct {
	ID        string `json:"trace_id"`
	Name      string `json:"name"`
	Status    string `json:"status,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Duration  int64  `json:"duration_ms,omitempty"`
	Model     string `json:"model,omitempty"`
	TotalTokens int   `json:"total_tokens,omitempty"`
	PromptTokens int  `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	Cost      float64 `json:"cost,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Workflow  string `json:"workflow,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Span struct {
	ID        string `json:"span_id"`
	TraceID   string `json:"trace_id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Input     any    `json:"input,omitempty"`
	Output    any    `json:"output,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Duration  int64  `json:"duration_ms"`
	Model     string `json:"model,omitempty"`
	ToolName  string `json:"tool_name,omitempty"`
	Status    string `json:"status,omitempty"`
	Error     string `json:"error,omitempty"`
	Children  []Span `json:"children,omitempty"`
}

type ListTracesParams struct {
	Workflow string
	Status   string
	Limit    int
	Offset   int
	From     string
	To       string
}

func (c *Client) ListTraces(params ListTracesParams) ([]Trace, error) {
	return nil, nil
}

func (c *Client) GetTrace(id string) (*Span, error) {
	return nil, nil
}
