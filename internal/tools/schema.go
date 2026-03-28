package tools

// FieldType 은 Tool 입력 필드의 데이터 타입이다.
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeObject  FieldType = "object"
	FieldTypeArray   FieldType = "array"
)

// FieldSchema 는 Tool 입력 스키마의 단일 필드를 기술한다.
type FieldSchema struct {
	Name        string    // 필드 키 이름 (tool input map의 key와 일치해야 함)
	Type        FieldType // 데이터 타입
	Description string    // LLM에게 전달할 필드 설명
	Required    bool      // true 이면 input map에 반드시 존재해야 함
}

// Schema 는 Tool 입력 전체의 구조를 기술한다.
// InputSchema() 의 반환 타입이며, ToolRouter 의 input validation 과
// Phase 3 LLM system prompt 직렬화에 함께 사용된다.
type Schema struct {
	Fields []FieldSchema
}
