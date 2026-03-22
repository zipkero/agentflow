package agent

import (
	"crypto/rand"
	"fmt"
)

const FixedSessionID = "session-dev"

// NewRequestID 는 crypto/rand 기반 UUID v4 형식의 요청 식별자를 생성한다.
func NewRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("request ID 생성 실패: %v", err))
	}
	// UUID v4: version bits(4), variant bits(2) 설정
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
