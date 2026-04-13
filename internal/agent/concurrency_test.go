package agent

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/zipkero/agent-runtime/internal/tools"
	"github.com/zipkero/agent-runtime/internal/types"
)

// slowTool 은 Execute 호출 시 지정된 시간만큼 대기하는 테스트용 tool 이다.
// context 가 취소되면 즉시 ctx.Err() 를 반환한다.
type slowTool struct {
	delay time.Duration
}

func (s *slowTool) Name() string        { return "slow_tool" }
func (s *slowTool) Description() string { return "sleeps for a configured duration" }
func (s *slowTool) InputSchema() tools.Schema {
	return tools.Schema{}
}
func (s *slowTool) Execute(ctx context.Context, _ map[string]any) (types.ToolResult, error) {
	select {
	case <-time.After(s.delay):
		return types.ToolResult{Output: "done"}, nil
	case <-ctx.Done():
		return types.ToolResult{}, ctx.Err()
	}
}

func newRouterWithSlowTool(delay time.Duration) *tools.ToolRouter {
	reg := tools.NewInMemoryToolRegistry()
	reg.Register(&slowTool{delay: delay})
	return tools.NewToolRouter(reg, slog.Default())
}

func TestToolTimeout_ExceedsDeadline(t *testing.T) {
	// 부모 context 에 매우 짧은 deadline 을 설정해 timeout 을 유발한다.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	router := newRouterWithSlowTool(5 * time.Second)
	_, err := router.Route(ctx, types.PlanResult{
		ActionType: types.ActionToolCall,
		ToolName:   "slow_tool",
		ToolInput:  map[string]any{},
	})

	if err == nil {
		t.Fatal("timeout 에러가 반환되어야 한다")
	}
	agentErr, ok := err.(*types.AgentError)
	if !ok {
		t.Fatalf("*types.AgentError 타입이어야 한다, got %T", err)
	}
	if agentErr.Kind != types.ErrToolTimeout {
		t.Errorf("에러 유형 불일치: got %q, want %q", agentErr.Kind, types.ErrToolTimeout)
	}
	if !agentErr.Retryable {
		t.Error("tool_timeout 은 retryable 이어야 한다")
	}
}

func TestToolTimeout_CompletesWithinDeadline(t *testing.T) {
	router := newRouterWithSlowTool(10 * time.Millisecond)

	result, err := router.Route(context.Background(), types.PlanResult{
		ActionType: types.ActionToolCall,
		ToolName:   "slow_tool",
		ToolInput:  map[string]any{},
	})

	if err != nil {
		t.Fatalf("에러 없이 완료되어야 한다: %v", err)
	}
	if result.Output != "done" {
		t.Errorf("결과 불일치: got %q, want %q", result.Output, "done")
	}
}

func TestToolTimeout_ParentCancelPropagates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	router := newRouterWithSlowTool(5 * time.Second)

	// 부모 context 를 50ms 후 취소한다.
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := router.Route(ctx, types.PlanResult{
		ActionType: types.ActionToolCall,
		ToolName:   "slow_tool",
		ToolInput:  map[string]any{},
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("부모 context 취소 시 에러가 반환되어야 한다")
	}
	// 5초 대기하지 않고 빠르게 반환되었는지 확인
	if elapsed > 1*time.Second {
		t.Errorf("부모 context 취소가 전파되지 않음: elapsed %v", elapsed)
	}
}
