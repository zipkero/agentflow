package agent

import "fmt"

// ErrorKind 는 AgentError 의 세부 유형이다.
type ErrorKind string

const (
	// ErrToolNotFound 는 registry 에 등록되지 않은 tool 을 호출한 경우다. fatal.
	ErrToolNotFound ErrorKind = "tool_not_found"
	// ErrInputValidationFailed 는 tool input 이 schema 를 만족하지 않는 경우다. fatal.
	ErrInputValidationFailed ErrorKind = "input_validation_failed"
	// ErrToolExecutionFailed 는 tool 실행 중 오류가 발생한 경우다. retryable.
	ErrToolExecutionFailed ErrorKind = "tool_execution_failed"
	// ErrLLMParseFailed 는 LLM 응답을 파싱할 수 없는 경우다. retryable.
	ErrLLMParseFailed ErrorKind = "llm_parse_error"
)

// AgentError 는 agent loop 내에서 발생하는 구조화된 에러 타입이다.
// Retryable 이 true 이면 loop 에서 재시도 가능, false 이면 즉시 종료.
type AgentError struct {
	Kind      ErrorKind
	Retryable bool
	Msg       string
}

func (e *AgentError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Kind, e.Msg)
}

func newFatalError(kind ErrorKind, msg string) *AgentError {
	return &AgentError{Kind: kind, Retryable: false, Msg: msg}
}

func newRetryableError(kind ErrorKind, msg string) *AgentError {
	return &AgentError{Kind: kind, Retryable: true, Msg: msg}
}

// NewToolNotFoundError 는 tool_not_found fatal 에러를 생성한다.
func NewToolNotFoundError(toolName string) *AgentError {
	return newFatalError(ErrToolNotFound, fmt.Sprintf("tool %q not found in registry", toolName))
}

// NewInputValidationError 는 input_validation_failed fatal 에러를 생성한다.
func NewInputValidationError(msg string) *AgentError {
	return newFatalError(ErrInputValidationFailed, msg)
}

// NewToolExecutionError 는 tool_execution_failed retryable 에러를 생성한다.
func NewToolExecutionError(toolName string, cause error) *AgentError {
	return newRetryableError(ErrToolExecutionFailed, fmt.Sprintf("tool %q execution failed: %v", toolName, cause))
}

// NewLLMParseError 는 llm_parse_error retryable 에러를 생성한다.
func NewLLMParseError(cause error) *AgentError {
	return newRetryableError(ErrLLMParseFailed, fmt.Sprintf("failed to parse LLM response: %v", cause))
}
