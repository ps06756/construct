package conv

import (
	"fmt"

	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/furisto/construct/backend/memory"
	"github.com/furisto/construct/backend/memory/schema/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertMemoryMessageToProto(m *memory.Message) (*v1.Message, error) {
	if m == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return &v1.Message{
		Metadata: &v1.MessageMetadata{
			Id:        m.ID.String(),
			CreatedAt: timestamppb.New(m.CreateTime),
			UpdatedAt: timestamppb.New(m.UpdateTime),
			TaskId:    m.TaskID.String(),
			AgentId:   ConvertUUIDPtrToStringPtr(m.AgentID),
			ModelId:   ConvertUUIDPtrToStringPtr(m.ModelID),
			Role:      convertRole(m.Source),
		},
		Spec: &v1.MessageSpec{
			Content: []*v1.MessagePart{
				{
					Data: &v1.MessagePart_Text_{
						Text: &v1.MessagePart_Text{
							Content: convertContent(m.Content),
						},
					},
				},
			},
		},
		Status: &v1.MessageStatus{
			Usage: convertUsage(m.Usage),
		},
	}, nil
}

func convertRole(role types.MessageSource) v1.MessageRole {
	switch role {
	case types.MessageSourceUser:
		return v1.MessageRole_MESSAGE_ROLE_USER
	case types.MessageSourceAssistant:
		return v1.MessageRole_MESSAGE_ROLE_ASSISTANT
	default:
		return v1.MessageRole_MESSAGE_ROLE_UNSPECIFIED
	}
}

func convertUsage(usage *types.MessageUsage) *v1.MessageUsage {
	if usage == nil {
		return nil
	}

	return &v1.MessageUsage{
		InputTokens:      usage.InputTokens,
		OutputTokens:     usage.OutputTokens,
		CacheWriteTokens: usage.CacheWriteTokens,
		CacheReadTokens:  usage.CacheReadTokens,
		Cost:             usage.Cost,
	}
}

func convertContent(content *types.MessageContent) string {
	if content == nil || len(content.Blocks) == 0 {
		return ""
	}

	for _, block := range content.Blocks {
		if block.Kind == types.MessageBlockKindText {
			return block.Payload
		}
	}
	return ""
}

func ConvertProtoMessageToMemory(m *v1.Message) (*memory.Message, error) {
	if m == nil {
		return nil, fmt.Errorf("message is nil")
	}

	messageID, err := ConvertStringToUUID(m.Metadata.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID: %w", err)
	}

	taskID, err := ConvertStringToUUID(m.Metadata.TaskId)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	agentID, err := ConvertStringPtrToUUID(m.Metadata.AgentId)
	if err != nil {
		return nil, fmt.Errorf("invalid agent ID: %w", err)
	}

	modelID, err := ConvertStringPtrToUUID(m.Metadata.ModelId)
	if err != nil {
		return nil, fmt.Errorf("invalid model ID: %w", err)
	}

	source, err := ConvertProtoRoleToMemoryRole(m.Metadata.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid message role: %w", err)
	}

	return &memory.Message{
		ID:         messageID,
		CreateTime: m.Metadata.CreatedAt.AsTime(),
		UpdateTime: m.Metadata.UpdatedAt.AsTime(),
		TaskID:     taskID,
		AgentID:    agentID,
		ModelID:    modelID,
		Source:     source,
		Content:    ConvertProtoContentToMemory(m.Spec.Content),
		Usage:      ConvertProtoUsageToMemoryUsage(m.Status.Usage),
	}, nil
}

func ConvertProtoRoleToMemoryRole(role v1.MessageRole) (types.MessageSource, error) {
	switch role {
	case v1.MessageRole_MESSAGE_ROLE_USER:
		return types.MessageSourceUser, nil
	case v1.MessageRole_MESSAGE_ROLE_ASSISTANT:
		return types.MessageSourceAssistant, nil
	default:
		return "", fmt.Errorf("invalid message role: %v", role)
	}
}

func ConvertProtoContentToMemory(content []*v1.MessagePart) *types.MessageContent {
	if len(content) == 0 {
		return nil
	}

	blocks := make([]types.MessageBlock, 0, len(content))
	for _, part := range content {
		if part.Data == nil {
			continue
		}

		switch messagePart := part.Data.(type) {
		case *v1.MessagePart_Text_:
			blocks = append(blocks, types.MessageBlock{
				Kind:    types.MessageBlockKindText,
				Payload: messagePart.Text.Content,
			})
		}
	}

	return &types.MessageContent{
		Blocks: blocks,
	}
}

func ConvertProtoUsageToMemoryUsage(usage *v1.MessageUsage) *types.MessageUsage {
	if usage == nil {
		return nil
	}

	return &types.MessageUsage{
		InputTokens:      usage.InputTokens,
		OutputTokens:     usage.OutputTokens,
		CacheWriteTokens: usage.CacheWriteTokens,
		CacheReadTokens:  usage.CacheReadTokens,
		Cost:             usage.Cost,
	}
}
