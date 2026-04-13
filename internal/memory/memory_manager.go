package memory

import (
	"context"

	"github.com/zipkero/agent-runtime/internal/state"
	"github.com/zipkero/agent-runtime/internal/types"
)

// MemoryManager 는 SessionRepository 와 MemoryRepository 를 단일 인터페이스로 캡슐화하는 파사드다.
// Runtime 은 이 인터페이스만 의존하며, 구체 저장소를 직접 알지 않는다.
type MemoryManager interface {
	// LoadSession 은 sessionID 에 해당하는 SessionState 를 반환한다.
	// 존재하지 않으면 빈 SessionState 와 nil error 를 반환한다.
	LoadSession(ctx context.Context, sessionID string) (state.SessionState, error)

	// SaveSession 은 sessionID 에 SessionState 를 저장한다.
	SaveSession(ctx context.Context, sessionID string, s state.SessionState) error

	// SaveMemory 는 Long-term Memory 레코드를 저장한다.
	SaveMemory(ctx context.Context, memory types.Memory) error

	// LoadRelevantMemory 는 userInput 을 기반으로 관련 Long-term Memory 를 조회한다.
	// 현재 구현은 userInput 에서 태그를 추출해 OR 조건으로 검색한다.
	LoadRelevantMemory(ctx context.Context, userInput string) ([]types.Memory, error)
}
