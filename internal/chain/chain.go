package chain

import (
	"context"
	"fmt"
	"sync"
)

// Step 定义链式调用中的单个步骤
type Step func(ctx context.Context, input interface{}) (interface{}, error)

// Chain 链式调用结构
type Chain struct {
	steps []Step
	mu    sync.RWMutex
}

// NewChain 创建新的链式调用
func NewChain() *Chain {
	return &Chain{
		steps: make([]Step, 0),
	}
}

// AddStep 添加步骤到链中
func (c *Chain) AddStep(step Step) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.steps = append(c.steps, step)
}

// Run 执行链式调用
func (c *Chain) Run(ctx context.Context, input interface{}) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := input
	var err error
	
	for i, step := range c.steps {
		result, err = step(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("step %d failed: %w", i, err)
		}
	}
	
	return result, nil
}

// RunString 执行链式调用（字符串输入输出）
func (c *Chain) RunString(ctx context.Context, input string) (string, error) {
	result, err := c.Run(ctx, input)
	if err != nil {
		return "", err
	}
	
	if str, ok := result.(string); ok {
		return str, nil
	}
	
	return fmt.Sprintf("%v", result), nil
}

// GetStepCount 获取步骤数量
func (c *Chain) GetStepCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.steps)
}

// Clear 清空所有步骤
func (c *Chain) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.steps = make([]Step, 0)
} 