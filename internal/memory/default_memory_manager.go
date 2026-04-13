package memory

import (
	"context"
	"strings"

	"github.com/zipkero/agent-runtime/internal/state"
	"github.com/zipkero/agent-runtime/internal/types"
)

const defaultRelevantMemoryLimit = 10

// DefaultMemoryManager 는 SessionRepository 와 MemoryRepository 를 주입받아
// MemoryManager 인터페이스를 구현하는 기본 구조체다.
type DefaultMemoryManager struct {
	sessions state.SessionRepository
	memories MemoryRepository
}

// NewDefaultMemoryManager 는 DefaultMemoryManager 를 생성한다.
func NewDefaultMemoryManager(sessions state.SessionRepository, memories MemoryRepository) *DefaultMemoryManager {
	return &DefaultMemoryManager{
		sessions: sessions,
		memories: memories,
	}
}

func (m *DefaultMemoryManager) LoadSession(ctx context.Context, sessionID string) (state.SessionState, error) {
	return m.sessions.Load(ctx, sessionID)
}

func (m *DefaultMemoryManager) SaveSession(ctx context.Context, sessionID string, s state.SessionState) error {
	return m.sessions.Save(ctx, sessionID, s)
}

func (m *DefaultMemoryManager) SaveMemory(ctx context.Context, memory types.Memory) error {
	return m.memories.Save(ctx, memory)
}

func (m *DefaultMemoryManager) LoadRelevantMemory(ctx context.Context, userInput string) ([]types.Memory, error) {
	tags := extractTags(userInput)
	if len(tags) == 0 {
		return nil, nil
	}
	return m.memories.LoadByTags(ctx, tags, defaultRelevantMemoryLimit)
}

// extractTags 는 userInput 을 공백으로 분리한 뒤 길이 2 이하의 토큰을 제외하여 태그 목록을 반환한다.
func extractTags(input string) []string {
	words := strings.Fields(input)
	tags := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.ToLower(w)
		if len(w) > 2 {
			tags = append(tags, w)
		}
	}
	return tags
}
