package rag

import (
	"context"
	"fmt"
	"strings"
)

// Document 文档结构
type Document struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
	Score    float64           `json:"score,omitempty"`
}

// Retriever 检索器接口
type Retriever interface {
	// Retrieve 检索相关文档
	Retrieve(ctx context.Context, query string, limit int) ([]Document, error)
	
	// AddDocument 添加文档
	AddDocument(ctx context.Context, doc *Document) error
	
	// RemoveDocument 删除文档
	RemoveDocument(ctx context.Context, id string) error
}

// SimpleRetriever 简单检索器（基于关键词匹配）
type SimpleRetriever struct {
	documents map[string]*Document
}

// NewSimpleRetriever 创建简单检索器
func NewSimpleRetriever() *SimpleRetriever {
	return &SimpleRetriever{
		documents: make(map[string]*Document),
	}
}

// Retrieve 实现检索接口
func (r *SimpleRetriever) Retrieve(ctx context.Context, query string, limit int) ([]Document, error) {
	if limit <= 0 {
		limit = 5
	}
	
	query = strings.ToLower(query)
	queryWords := strings.Fields(query)
	
	var results []Document
	scores := make(map[string]float64)
	
	// 简单的关键词匹配评分
	for id, doc := range r.documents {
		content := strings.ToLower(doc.Content)
		score := 0.0
		
		for _, word := range queryWords {
			if strings.Contains(content, word) {
				score += 1.0
			}
		}
		
		if score > 0 {
			scores[id] = score
		}
	}
	
	// 按分数排序并返回结果
	for id, score := range scores {
		doc := *r.documents[id]
		doc.Score = score
		results = append(results, doc)
		
		if len(results) >= limit {
			break
		}
	}
	
	return results, nil
}

// AddDocument 添加文档
func (r *SimpleRetriever) AddDocument(ctx context.Context, doc *Document) error {
	if doc == nil {
		return fmt.Errorf("document cannot be nil")
	}
	
	if doc.ID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}
	
	r.documents[doc.ID] = doc
	return nil
}

// RemoveDocument 删除文档
func (r *SimpleRetriever) RemoveDocument(ctx context.Context, id string) error {
	if _, exists := r.documents[id]; !exists {
		return fmt.Errorf("document '%s' not found", id)
	}
	
	delete(r.documents, id)
	return nil
}

// RAGEngine RAG 引擎
type RAGEngine struct {
	retriever Retriever
}

// NewRAGEngine 创建 RAG 引擎
func NewRAGEngine(retriever Retriever) *RAGEngine {
	return &RAGEngine{
		retriever: retriever,
	}
}

// Query 执行 RAG 查询
func (e *RAGEngine) Query(ctx context.Context, query string, limit int) (string, error) {
	// 检索相关文档
	docs, err := e.retriever.Retrieve(ctx, query, limit)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}
	
	if len(docs) == 0 {
		return "抱歉，没有找到相关的文档信息。", nil
	}
	
	// 构建上下文
	context := e.buildContext(docs)
	
	// 构建增强的查询
	enhancedQuery := fmt.Sprintf("基于以下上下文信息回答问题：\n\n上下文：\n%s\n\n问题：%s", context, query)
	
	return enhancedQuery, nil
}

// buildContext 构建上下文
func (e *RAGEngine) buildContext(docs []Document) string {
	var context strings.Builder
	
	for i, doc := range docs {
		context.WriteString(fmt.Sprintf("文档 %d (相关度: %.2f):\n", i+1, doc.Score))
		context.WriteString(doc.Content)
		context.WriteString("\n\n")
	}
	
	return context.String()
}

// Retrieve 检索步骤（链式调用）
func Retrieve(ctx context.Context, input interface{}) (interface{}, error) {
	if str, ok := input.(string); ok {
		// 这里应该使用实际的检索器
		// 为了演示，返回一个简单的增强查询
		enhancedQuery := fmt.Sprintf("基于知识库检索，回答以下问题：\n\n%s", str)
		return enhancedQuery, nil
	}
	
	return input, nil
}

// AddDocumentToRAG 添加文档到 RAG 系统
func AddDocumentToRAG(ctx context.Context, input interface{}) (interface{}, error) {
	// 这里可以作为链式调用中的一个步骤
	// 将文档添加到 RAG 系统中
	return input, nil
} 