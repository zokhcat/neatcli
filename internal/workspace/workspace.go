package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DirName      = ".neatlogs"
	PromptsDir   = "prompts"
	ToolsDir     = "tools"
	WorkspaceYML = "workspace.yaml"
)

type PromptFile struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Labels   []string `yaml:"labels,omitempty"`
	Content  string   `yaml:"content,omitempty"`
	Messages []Message `yaml:"messages,omitempty"`
}

type Message struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

type ToolFile struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Schema      string `yaml:"schema,omitempty"`
}

type Workspace struct {
	root string
}

func Open() (*Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(cwd, DirName)
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("not a neatcli workspace (no %s/ directory): %w", DirName, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}
	return &Workspace{root: root}, nil
}

func Init() (*Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(cwd, DirName)

	dirs := []string{root, filepath.Join(root, PromptsDir), filepath.Join(root, ToolsDir)}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", d, err)
		}
	}

	return &Workspace{root: root}, nil
}

func (w *Workspace) Root() string { return w.root }

func (w *Workspace) PromptsDir() string { return filepath.Join(w.root, PromptsDir) }

func (w *Workspace) ToolsDir() string { return filepath.Join(w.root, ToolsDir) }

func (w *Workspace) ListPromptFiles() ([]string, error) {
	dir := w.PromptsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".yaml" {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (w *Workspace) ReadPromptFile(name string) (*PromptFile, error) {
	path := filepath.Join(w.PromptsDir(), name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf PromptFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &pf, nil
}

func (w *Workspace) WritePromptFile(pf *PromptFile) error {
	path := filepath.Join(w.PromptsDir(), pf.Name+".yaml")
	data, err := yaml.Marshal(pf)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (w *Workspace) ListToolFiles() ([]string, error) {
	dir := w.ToolsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".yaml" {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (w *Workspace) ReadToolFile(name string) (*ToolFile, error) {
	path := filepath.Join(w.ToolsDir(), name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tf ToolFile
	if err := yaml.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &tf, nil
}

func (w *Workspace) WriteToolFile(tf *ToolFile) error {
	path := filepath.Join(w.ToolsDir(), tf.Name+".yaml")
	data, err := yaml.Marshal(tf)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
