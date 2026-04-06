package types

// ActionType 은 Planner 가 결정할 수 있는 행동 유형이다.
type ActionType string

const (
	// ActionToolCall 은 특정 Tool 을 호출하도록 지시한다.
	ActionToolCall ActionType = "tool_call"
	// ActionRespondDirectly 는 Tool 없이 바로 응답을 생성한다.
	ActionRespondDirectly ActionType = "respond_directly"
	// ActionFinish 는 loop 를 종료한다.
	ActionFinish ActionType = "finish"
	// ActionSummarize 는 지금까지의 ToolResults 를 요약해 응답을 생성한다.
	// Executor 를 호출하지 않고 respond_directly 와 동일하게 loop 를 종료한다.
	ActionSummarize ActionType = "summarize"
	// ActionAskUser 는 사용자에게 추가 입력을 요청한다.
	// CLI 환경에서는 FinalAnswer 에 질문 문자열을 채우고 loop 를 즉시 종료한다.
	// HTTP API 환경에서의 비동기 대기 메커니즘은 Phase 7 에서 구현한다.
	ActionAskUser ActionType = "ask_user"
)

// PlanResult 는 Planner 가 내린 결정을 담는 구조체다.
// ActionType 이 tool_call 일 때 ToolName 과 ToolInput 이 유효하다.
// ActionType 이 respond_directly / summarize / ask_user 일 때 Reasoning 이 FinalAnswer 로 사용된다.
type PlanResult struct {
	// ActionType 은 이번 step 에서 취할 행동 유형이다.
	ActionType ActionType `json:"action_type"`
	// ToolName 은 ActionType 이 tool_call 일 때 호출할 tool 이름이다.
	ToolName string `json:"tool_name,omitempty"`
	// ToolInput 은 tool 에 전달할 인자다.
	ToolInput map[string]any `json:"tool_input,omitempty"`
	// Reasoning 은 LLM 이 이 결정을 내린 전체 추론 과정이다.
	// respond_directly / summarize / ask_user 일 때 FinalAnswer 로 사용된다.
	Reasoning string `json:"reasoning"`
	// ReasoningSummary 는 Reasoning 의 한 줄 요약이다.
	// 로그·디버깅용으로, prompt_builder 가 다음 step system prompt 에 포함한다.
	ReasoningSummary string `json:"reasoning_summary,omitempty"`
	// Confidence 는 LLM 이 스스로 평가한 결정 신뢰도다 (0.0 ~ 1.0).
	// RetryPolicy 와 Reflector 가 재시도 여부 판단에 활용한다.
	Confidence float64 `json:"confidence,omitempty"`
	// NextGoal 은 다음 step 에서 달성해야 할 목표를 LLM 이 미리 기술한 것이다.
	// prompt_builder 가 다음 step system prompt 에 포함해 연속성을 유지한다.
	NextGoal string `json:"next_goal,omitempty"`
}
