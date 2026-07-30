package api

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema,omitempty"`
	Source      string `json:"source,omitempty"`
}

func (c *Client) ListTools() ([]Tool, error) {
	return nil, nil
}

func (c *Client) GetTool(name string) (*Tool, error) {
	return nil, nil
}

func (c *Client) UpdateToolDescription(name, description string) error {
	return nil
}
