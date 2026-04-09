package state

import "time"

// SessionState 는 세션 범위의 상태를 담는다.
// 여러 번의 Run() 호출을 넘어 지속되는 데이터를 관리한다.
type SessionState struct {
	// SessionID 는 세션을 식별하는 고유 ID다.
	SessionID string
	// RecentContext 는 직전 대화 교환들의 요약 목록이다.
	// 오래된 순서부터 정렬되며, 새 교환이 추가될 때 앞에서부터 제거된다.
	RecentContext []string
	// ActiveGoal 은 현재 세션에서 사용자가 달성하려는 목표다.
	// 빈 문자열이면 명시적 목표 없음.
	ActiveGoal string
	// LastUpdated 는 이 세션 상태가 마지막으로 저장된 시각이다.
	LastUpdated time.Time
}
