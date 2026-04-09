package state

import "context"

// SessionRepository 는 SessionState의 저장과 조회를 추상화한다.
// in-memory, Redis 등 다양한 백엔드로 교체할 수 있도록 인터페이스로 분리한다.
type SessionRepository interface {
	// Load 는 sessionID에 해당하는 SessionState를 반환한다.
	// 존재하지 않으면 빈 SessionState와 nil error를 반환한다.
	Load(ctx context.Context, sessionID string) (SessionState, error)
	// Save 는 sessionID에 SessionState를 저장한다.
	Save(ctx context.Context, sessionID string, state SessionState) error
}
