package api

import (
	"fmt"
	"net/url"
	"strconv"
)

type Prompt struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Version   int                 `json:"version"`
	Type      string              `json:"type"`
	Content   string              `json:"content"`
	Messages  []map[string]string `json:"messages,omitempty"`
	Labels    []string            `json:"labels"`
	Tags      []string            `json:"tags,omitempty"`
	Config    map[string]any      `json:"config,omitempty"`
	CreatedAt string              `json:"createdAt,omitempty"`
	UpdatedAt string              `json:"updatedAt,omitempty"`
}

type ListPromptsResponse struct {
	Items []Prompt `json:"items"`
	Total int      `json:"total,omitempty"`
}

type CreatePromptRequest struct {
	Name          string              `json:"name"`
	Type          string              `json:"type"`
	Content       string              `json:"content,omitempty"`
	Messages      []map[string]string `json:"messages,omitempty"`
	Labels        []string            `json:"labels"`
	Tags          []string            `json:"tags,omitempty"`
	Config        map[string]any      `json:"config,omitempty"`
	CommitMessage string              `json:"commit_message,omitempty"`
}

type SaveVersionRequest struct {
	PromptName    string              `json:"promptName"`
	Content       string              `json:"content,omitempty"`
	Messages      []map[string]string `json:"messages,omitempty"`
	Labels        []string            `json:"labels,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
	Config        map[string]any      `json:"config,omitempty"`
	CommitMessage string              `json:"commitMessage,omitempty"`
}

type LabelRequest struct {
	Label string `json:"label"`
}

func (c *Client) ListPrompts(name, label string, limit, offset int) (*ListPromptsResponse, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	if name != "" {
		params.Set("name", name)
	}
	if label != "" {
		params.Set("label", label)
	}

	var resp ListPromptsResponse
	if err := c.requestJSON("GET", "/api/managed-prompts", params, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetPrompt(name, label string, version int) (*Prompt, error) {
	if label != "" {
		path := fmt.Sprintf("/api/v1/prompts/%s/fetch", url.PathEscape(name))
		params := url.Values{"label": {label}}
		var result map[string]any
		if err := c.requestJSON("GET", path, params, nil, &result); err != nil {
			return nil, err
		}
		p := normalizePrompt(result)
		return &p, nil
	}

	if version > 0 {
		resp, err := c.ListPrompts(name, "", 100, 0)
		if err != nil {
			return nil, err
		}
		for _, p := range resp.Items {
			if p.Version == version {
				return &p, nil
			}
		}
		return nil, fmt.Errorf("prompt %q version %d not found", name, version)
	}

	resp, err := c.ListPrompts(name, "", 1, 0)
	if err != nil {
		return nil, err
	}
	items := resp.Items
	if len(items) == 0 {
		return nil, fmt.Errorf("prompt %q not found", name)
	}

	latest := items[0]
	for _, p := range items {
		if p.CreatedAt > latest.CreatedAt {
			latest = p
		}
	}
	return &latest, nil
}

func (c *Client) CreatePrompt(req *CreatePromptRequest) (*Prompt, error) {
	var result map[string]any
	if err := c.requestJSON("POST", "/api/managed-prompts", nil, req, &result); err != nil {
		return nil, err
	}

	promptRaw, ok := result["prompt"]
	if !ok {
		promptRaw = result
	}
	if m, ok := promptRaw.(map[string]any); ok {
		p := normalizePrompt(m)
		return &p, nil
	}
	return nil, fmt.Errorf("unexpected create response shape")
}

func (c *Client) UpdatePromptLabels(name string, version int, newLabels []string) error {
	resp, err := c.ListPrompts(name, "", 100, 0)
	if err != nil {
		return err
	}

	var promptID string
	for _, p := range resp.Items {
		if p.Version == version {
			promptID = p.ID
			break
		}
	}
	if promptID == "" {
		return fmt.Errorf("prompt %q version %d not found", name, version)
	}

	for _, label := range newLabels {
		path := fmt.Sprintf("/api/managed-prompts/%s/labels", url.PathEscape(promptID))
		req := LabelRequest{Label: label}
		if _, err := c.request("POST", path, nil, &req); err != nil {
			return fmt.Errorf("add label %q: %w", label, err)
		}
	}
	return nil
}

func (c *Client) SaveAsVersion(req *SaveVersionRequest) (*Prompt, error) {
	var result map[string]any
	if err := c.requestJSON("POST", "/api/prompt-playground/save-as-version", nil, req, &result); err != nil {
		return nil, err
	}
	p := normalizePrompt(result)
	return &p, nil
}

func (c *Client) DeletePrompt(name string, version int) error {
	resp, err := c.ListPrompts(name, "", 100, 0)
	if err != nil {
		return err
	}

	var promptID string
	for _, p := range resp.Items {
		if p.Version == version {
			promptID = p.ID
			break
		}
	}
	if promptID == "" {
		return fmt.Errorf("prompt %q version %d not found", name, version)
	}

	path := fmt.Sprintf("/api/managed-prompts/%s", url.PathEscape(promptID))
	_, err = c.request("DELETE", path, nil, nil)
	return err
}

func normalizePrompt(raw map[string]any) Prompt {
	p := Prompt{
		ID:        getString(raw, "id"),
		Name:      getString(raw, "name"),
		Type:      getString(raw, "type"),
		Content:   getString(raw, "content"),
		CreatedAt: getString(raw, "createdAt", "created_at"),
		UpdatedAt: getString(raw, "updatedAt", "updated_at"),
	}

	if v, ok := raw["version"].(float64); ok {
		p.Version = int(v)
	}

	if labels, ok := raw["labels"].([]any); ok {
		for _, l := range labels {
			if s, ok := l.(string); ok {
				p.Labels = append(p.Labels, s)
			}
		}
	}

	if msgs, ok := raw["messages"].([]any); ok {
		for _, m := range msgs {
			if mMap, ok := m.(map[string]any); ok {
				msg := map[string]string{}
				if r, ok := mMap["role"].(string); ok {
					msg["role"] = r
				}
				if c, ok := mMap["content"].(string); ok {
					msg["content"] = c
				}
				p.Messages = append(p.Messages, msg)
			}
		}
	}

	if tags, ok := raw["tags"].([]any); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				p.Tags = append(p.Tags, s)
			}
		}
	}

	if cfg, ok := raw["config"].(map[string]any); ok {
		p.Config = cfg
	}

	return p
}

func getString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok {
			return v
		}
	}
	return ""
}
