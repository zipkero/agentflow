package tools

import (
	"context"
	"fmt"

	"agentflow/internal/agent"
	"agentflow/internal/planner"
	"agentflow/internal/state"
)

// ToolRouter 는 PlanResult 를 받아 registry 에서 tool 을 조회하고 실행한다.
// planner 와 tool 구현체 사이를 중재하며, 에러를 유형별로 분류해 반환한다.
type ToolRouter struct {
	registry ToolRegistry
}

func NewToolRouter(registry ToolRegistry) *ToolRouter {
	return &ToolRouter{registry: registry}
}

// Route 는 PlanResult 의 ToolName 으로 tool 을 조회하고 ToolInput 을 검증한 뒤 실행한다.
//
// 에러 유형:
//   - tool_not_found  : registry 에 없는 이름 → fatal
//   - input_validation_failed : required 필드 누락 또는 타입 불일치 → fatal
//   - tool_execution_failed   : Execute() 에서 error 반환 → retryable
func (r *ToolRouter) Route(ctx context.Context, plan planner.PlanResult) (state.ToolResult, error) {
	tool, err := r.registry.Get(plan.ToolName)
	if err != nil {
		return state.ToolResult{}, agent.NewToolNotFoundError(plan.ToolName)
	}

	if err := validateInput(tool.InputSchema(), plan.ToolInput); err != nil {
		return state.ToolResult{}, agent.NewInputValidationError(err.Error())
	}

	result, err := tool.Execute(ctx, plan.ToolInput)
	if err != nil {
		return state.ToolResult{}, agent.NewToolExecutionError(plan.ToolName, err)
	}

	return result, nil
}

// validateInput 은 schema 의 required 필드 존재 여부와 타입을 검증한다.
func validateInput(schema Schema, input map[string]any) error {
	for _, field := range schema.Fields {
		val, ok := input[field.Name]
		if !ok {
			if field.Required {
				return fmt.Errorf("required field %q is missing", field.Name)
			}
			continue
		}
		if err := checkType(field, val); err != nil {
			return err
		}
	}
	return nil
}

func checkType(field FieldSchema, val any) error {
	switch field.Type {
	case FieldTypeString:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("field %q must be string, got %T", field.Name, val)
		}
	case FieldTypeNumber:
		switch val.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
		default:
			return fmt.Errorf("field %q must be number, got %T", field.Name, val)
		}
	case FieldTypeBoolean:
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("field %q must be boolean, got %T", field.Name, val)
		}
	}
	return nil
}
