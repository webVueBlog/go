package prompt

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	_ "strings"
	"text/template"
)

// Template Prompt 模板
type Template struct {
	Name      string            `json:"name"`
	Content   string            `json:"content"`
	Version   string            `json:"version"`
	Variables []string          `json:"variables"`
	Metadata  map[string]string `json:"metadata"`
}

// PromptEngine Prompt 引擎
type PromptEngine struct {
	templates map[string]*Template
}

// NewPromptEngine 创建 Prompt 引擎
func NewPromptEngine() *PromptEngine {
	return &PromptEngine{
		templates: make(map[string]*Template),
	}
}

// AddTemplate 添加模板
func (e *PromptEngine) AddTemplate(tmpl *Template) error {
	if tmpl == nil {
		return fmt.Errorf("template cannot be nil")
	}

	if tmpl.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if tmpl.Content == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// 解析变量
	tmpl.Variables = e.extractVariables(tmpl.Content)

	e.templates[tmpl.Name] = tmpl
	return nil
}

// GetTemplate 获取模板
func (e *PromptEngine) GetTemplate(name string) (*Template, error) {
	tmpl, exists := e.templates[name]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found", name)
	}
	return tmpl, nil
}

// Render 渲染模板
func (e *PromptEngine) Render(name string, data map[string]interface{}) (string, error) {
	tmpl, err := e.GetTemplate(name)
	if err != nil {
		return "", err
	}

	// 使用 Go template 引擎
	t, err := template.New(tmpl.Name).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderWithVariables 渲染模板（使用变量替换）
func (e *PromptEngine) RenderWithVariables(name string, variables map[string]string) (string, error) {
	// 转换变量格式
	data := make(map[string]interface{})
	for k, v := range variables {
		data[k] = v
	}

	return e.Render(name, data)
}

// extractVariables 提取模板中的变量
func (e *PromptEngine) extractVariables(content string) []string {
	// 匹配 {{.variable}} 格式的变量
	re := regexp.MustCompile(`\{\{\.(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	variables := make([]string, 0)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			variable := match[1]
			if !seen[variable] {
				variables = append(variables, variable)
				seen[variable] = true
			}
		}
	}

	return variables
}

// ListTemplates 列出所有模板
func (e *PromptEngine) ListTemplates() []string {
	names := make([]string, 0, len(e.templates))
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}

// RemoveTemplate 删除模板
func (e *PromptEngine) RemoveTemplate(name string) error {
	if _, exists := e.templates[name]; !exists {
		return fmt.Errorf("template '%s' not found", name)
	}

	delete(e.templates, name)
	return nil
}

// BuildPrompt 构建 Prompt（链式调用步骤）
func BuildPrompt(ctx context.Context, input interface{}) (interface{}, error) {
	// 这里可以作为链式调用中的一个步骤
	// 根据输入构建合适的 Prompt
	if str, ok := input.(string); ok {
		// 简单的 Prompt 构建逻辑
		prompt := fmt.Sprintf("请回答以下问题：\n\n%s", str)
		return prompt, nil
	}

	return input, nil
}

// 预定义模板
var DefaultTemplates = map[string]*Template{
	"qa": {
		Name:    "qa",
		Content: "请回答以下问题：\n\n{{.question}}\n\n请提供详细、准确的答案。",
		Version: "1.0",
		Metadata: map[string]string{
			"type": "question-answer",
		},
	},
	"translation": {
		Name:    "translation",
		Content: "请将以下文本翻译成{{.target_language}}：\n\n{{.text}}\n\n请保持原文的意思和风格。",
		Version: "1.0",
		Metadata: map[string]string{
			"type": "translation",
		},
	},
	"summary": {
		Name:    "summary",
		Content: "请总结以下文本的主要内容：\n\n{{.text}}\n\n请提供简洁、准确的总结。",
		Version: "1.0",
		Metadata: map[string]string{
			"type": "summarization",
		},
	},
	"code_review": {
		Name:    "code_review",
		Content: "请对以下代码进行审查：\n\n```{{.language}}\n{{.code}}\n```\n\n请从代码质量、安全性、性能等方面进行评估，并提供改进建议。",
		Version: "1.0",
		Metadata: map[string]string{
			"type": "code-review",
		},
	},
}
