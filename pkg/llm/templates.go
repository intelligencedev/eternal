package llm

import (
	"fmt"
	"strings"
)

// Message represents a message for the completion API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PromptTemplate represents a template for generating string prompts.
type PromptTemplate struct {
	Template string
}

// ChatPromptTemplate represents a template for generating chat prompts.
type ChatPromptTemplate struct {
	Messages []Message
}

// GetSystemTemplate returns the system template.
func GetSystemTemplate(userPrompt string) ChatPromptTemplate {
	userPrompt = fmt.Sprintf("{%s}", userPrompt)
	template := NewChatPromptTemplate([]Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant and expert researcher and developer.",
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	})

	return *template
}

// NewChatPromptTemplate creates a new ChatPromptTemplate.
func NewChatPromptTemplate(messages []Message) *ChatPromptTemplate {
	return &ChatPromptTemplate{Messages: messages}
}

// Format formats the template with the provided variables.
func (pt *PromptTemplate) Format(vars map[string]string) string {
	result := pt.Template
	for k, v := range vars {
		placeholder := fmt.Sprintf("{%s}", k)
		result = strings.ReplaceAll(result, placeholder, v)
	}
	return result
}

// FormatMessages formats the chat messages with the provided variables.
func (cpt *ChatPromptTemplate) FormatMessages(vars map[string]string) []Message {
	var formattedMessages []Message
	for _, msg := range cpt.Messages {
		formattedContent := msg.Content
		for k, v := range vars {
			placeholder := fmt.Sprintf("{%s}", k)
			formattedContent = strings.ReplaceAll(formattedContent, placeholder, v)
		}
		formattedMessages = append(formattedMessages, Message{Role: msg.Role, Content: formattedContent})
	}
	return formattedMessages
}
