package internal

import (
	"context"

	"github.com/anthropics/claude-code-sdk-go/model"
)

// Client orchestrates communication with the Claude CLI.
type Client struct {
	Transport Transport
}

// Query connects to the transport and streams parsed messages.
func (c *Client) Query(ctx context.Context, prompt string, opts *model.Options) (<-chan model.Message, error) {
	if sp, ok := c.Transport.(*SubprocessCLITransport); ok {
		sp.Prompt = prompt
		sp.Options = opts
	}

	if err := c.Transport.Connect(ctx); err != nil {
		return nil, err
	}

	rawCh, err := c.Transport.ReceiveMessages(ctx)
	if err != nil {
		return nil, err
	}

	out := make(chan model.Message)
	go func() {
		defer close(out)
		defer c.Transport.Disconnect()
		for data := range rawCh {
			if m := parseMessage(data); m != nil {
				out <- m
			}
		}
	}()
	return out, nil
}

func parseMessage(data map[string]any) model.Message {
	t, _ := data["type"].(string)
	switch t {
	case "user":
		if msg, ok := data["message"].(map[string]any); ok {
			if content, ok2 := msg["content"].(string); ok2 {
				return model.UserMessage{Content: content}
			}
		}
	case "assistant":
		var blocks []model.ContentBlock
		if msg, ok := data["message"].(map[string]any); ok {
			if arr, ok2 := msg["content"].([]any); ok2 {
				for _, raw := range arr {
					blockMap, ok3 := raw.(map[string]any)
					if !ok3 {
						continue
					}
					bType, _ := blockMap["type"].(string)
					switch bType {
					case "text":
						if txt, ok := blockMap["text"].(string); ok {
							blocks = append(blocks, model.TextBlock{Text: txt})
						}
					case "tool_use":
						bu := model.ToolUseBlock{}
						if v, ok := blockMap["id"].(string); ok {
							bu.ID = v
						}
						if v, ok := blockMap["name"].(string); ok {
							bu.Name = v
						}
						if v, ok := blockMap["input"].(map[string]any); ok {
							bu.Input = v
						}
						blocks = append(blocks, bu)
					case "tool_result":
						br := model.ToolResultBlock{}
						if v, ok := blockMap["tool_use_id"].(string); ok {
							br.ToolUseID = v
						}
						if v, ok := blockMap["content"]; ok {
							br.Content = v
						}
						if v, ok := blockMap["is_error"].(bool); ok {
							br.IsError = v
						}
						blocks = append(blocks, br)
					}
				}
			}
		}
		return model.AssistantMessage{Content: blocks}
	case "system":
		subtype, _ := data["subtype"].(string)
		return model.SystemMessage{Subtype: subtype, Data: data}
	case "result":
		msg := model.ResultMessage{}
		if v, ok := data["subtype"].(string); ok {
			msg.Subtype = v
		}
		if v, ok := data["cost_usd"].(float64); ok {
			msg.CostUSD = v
		}
		if v, ok := data["duration_ms"].(float64); ok {
			msg.DurationMS = int(v)
		}
		if v, ok := data["duration_api_ms"].(float64); ok {
			msg.DurationAPIMS = int(v)
		}
		if v, ok := data["is_error"].(bool); ok {
			msg.IsError = v
		}
		if v, ok := data["num_turns"].(float64); ok {
			msg.NumTurns = int(v)
		}
		if v, ok := data["session_id"].(string); ok {
			msg.SessionID = v
		}
		if v, ok := data["total_cost"].(float64); ok {
			msg.TotalCostUSD = v
		}
		if v, ok := data["usage"].(map[string]any); ok {
			msg.Usage = v
		}
		if v, ok := data["result"].(string); ok {
			msg.Result = v
		}
		return msg
	}
	return nil
}
