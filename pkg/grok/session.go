package grok

import (
	"context"

	"github.com/google/uuid"
)

func GenerateSessionID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}

func (c *GrokClient) ContinueConversation(prompt string) (*GrokResult, error) {
	return c.ContinueConversationCtx(context.Background(), prompt)
}

func (c *GrokClient) ContinueConversationCtx(ctx context.Context, prompt string) (*GrokResult, error) {
	return c.RunPromptCtx(ctx, prompt, &RunOptions{
		Format:   JSONOutput,
		Continue: true,
	})
}

func (c *GrokClient) ResumeConversation(prompt, sessionID string) (*GrokResult, error) {
	return c.ResumeConversationCtx(context.Background(), prompt, sessionID)
}

func (c *GrokClient) ResumeConversationCtx(ctx context.Context, prompt, sessionID string) (*GrokResult, error) {
	return c.RunPromptCtx(ctx, prompt, &RunOptions{
		Format:   JSONOutput,
		ResumeID: sessionID,
	})
}
