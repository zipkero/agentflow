package planner

// PlanResultSchema 는 LLM 이 반환해야 할 JSON 의 구조를 JSON Schema 형식으로 기술한 문자열이다.
// system prompt 에 삽입해 LLM 이 형식을 어기는 것을 예방한다.
const PlanResultSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["action_type", "reasoning"],
  "properties": {
    "action_type": {
      "type": "string",
      "enum": ["tool_call", "respond_directly", "finish", "summarize", "ask_user"],
      "description": "이번 step 에서 취할 행동 유형"
    },
    "tool_name": {
      "type": "string",
      "description": "action_type 이 tool_call 일 때만 필수. 호출할 tool 의 이름"
    },
    "tool_input": {
      "type": "object",
      "description": "action_type 이 tool_call 일 때만 필수. tool 에 전달할 인자 (key-value)"
    },
    "reasoning": {
      "type": "string",
      "description": "이 결정을 내린 전체 추론 과정. respond_directly / summarize / ask_user 일 때는 이 값이 최종 응답 또는 질문 문자열로 사용됨"
    },
    "reasoning_summary": {
      "type": "string",
      "description": "reasoning 의 한 줄 요약. 로그 및 다음 step 프롬프트에 포함됨"
    },
    "confidence": {
      "type": "number",
      "minimum": 0.0,
      "maximum": 1.0,
      "description": "이 결정에 대한 신뢰도 (0.0 ~ 1.0)"
    },
    "next_goal": {
      "type": "string",
      "description": "다음 step 에서 달성해야 할 목표. 연속성 유지를 위해 다음 step 프롬프트에 포함됨"
    }
  },
  "if": {
    "properties": { "action_type": { "const": "tool_call" } },
    "required": ["action_type"]
  },
  "then": {
    "required": ["tool_name", "tool_input"]
  }
}`

// PlanResultSchemaExamples 는 각 action_type 별 예시 JSON 이다.
// JSON Schema 만으로는 LLM 이 의도를 오해하는 경우가 있어 구체적인 예시를 함께 제공한다.
const PlanResultSchemaExamples = `
## 예시 1 — tool_call (tool 을 호출해야 할 때)
{
  "action_type": "tool_call",
  "tool_name": "weather_mock",
  "tool_input": { "city": "Seoul" },
  "reasoning": "사용자가 서울 날씨를 물어봤으므로 weather_mock tool 을 호출한다.",
  "reasoning_summary": "서울 날씨 조회",
  "confidence": 0.95,
  "next_goal": "날씨 데이터를 받은 후 자연어로 응답한다."
}

## 예시 2 — respond_directly (tool 없이 바로 답할 수 있을 때)
{
  "action_type": "respond_directly",
  "reasoning": "서울의 현재 날씨는 맑고 기온은 22도입니다.",
  "reasoning_summary": "날씨 정보 응답",
  "confidence": 0.9
}

## 예시 3 — finish (목표가 완전히 달성됐을 때)
{
  "action_type": "finish",
  "reasoning": "모든 정보를 수집하고 사용자에게 전달 완료했다.",
  "reasoning_summary": "태스크 완료",
  "confidence": 1.0
}

## 예시 4 — summarize (여러 tool 결과를 종합해 답해야 할 때)
{
  "action_type": "summarize",
  "reasoning": "수집된 검색 결과를 바탕으로 전체 내용을 요약해 응답한다.",
  "reasoning_summary": "검색 결과 요약",
  "confidence": 0.85
}

## 예시 5 — ask_user (사용자에게 추가 정보가 필요할 때)
{
  "action_type": "ask_user",
  "reasoning": "어느 날짜의 날씨를 원하시나요?",
  "reasoning_summary": "날짜 정보 요청",
  "confidence": 0.8
}`

// PlanResultSchemaPrompt 는 system prompt 에 삽입할 준비가 된 스키마 + 예시 블록을 반환한다.
// LLMPlanner 의 prompt_builder 에서 이 함수를 호출해 시스템 프롬프트에 포함한다.
func PlanResultSchemaPrompt() string {
	return `## 응답 형식

반드시 아래 JSON Schema 를 따르는 JSON 객체만 반환하라. 설명, 마크다운 코드블록, 추가 텍스트 없이 JSON 만 출력하라.

### JSON Schema
` + PlanResultSchema + `

### action_type 선택 기준
- tool_call: 정보 수집이나 계산을 위해 tool 을 호출해야 할 때
- respond_directly: tool 없이 현재 정보만으로 충분히 답할 수 있을 때
- summarize: 여러 tool 실행 결과를 종합해 최종 응답을 만들어야 할 때
- ask_user: 사용자 의도가 불명확하거나 필수 정보가 누락됐을 때
- finish: 모든 목표가 달성돼 대화를 종료해야 할 때
` + PlanResultSchemaExamples
}
