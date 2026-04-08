package state

// SessionState 는 세션 범위의 상태를 담는다.
// 여러 번의 Run() 호출을 넘어 지속되는 데이터를 관리한다.
// Phase 4-2에서 RecentContext, ActiveGoal, LastUpdated 필드가 추가된다.
type SessionState struct {
	// SessionID 는 세션을 식별하는 고유 ID다.
	SessionID string
}
