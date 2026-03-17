# Scope — 범위 고정 문서

이 프로젝트가 다루는 것과 다루지 않는 것을 명확히 기술한다.
나중에 scope creep이 발생했을 때 이 문서를 기준으로 판단한다.

---

## 다루는 것 (In Scope)

### Agent 유형
- **QA형 Agent** — 질문에 답변하기 위해 Tool을 선택하고 실행하는 Agent
- **Search형 Agent** — 검색 결과를 수집, 필터, 요약하는 Agent
- **Planning형 Agent** — 복잡한 목표를 Task로 분해하고 순서대로 실행하는 Agent

### Runtime 핵심 구조
- Agent Loop (plan → execute → verify → finish)
- Planner 인터페이스 및 LLM 기반 구현
- Tool Registry / Tool Router / Tool 인터페이스
- AgentState / RequestState / SessionState / WorkingMemory / LongTermMemory 분리
- Verifier / RetryPolicy / FailureHandler
- Multi-Agent Orchestration (Task 분해 + 병렬 실행 + 결과 병합)

### 인프라
- Redis — Session 상태 저장
- Postgres — Long-term Memory 저장
- Docker Compose — 로컬 인프라 구동

### 실행 인터페이스
- Phase 1~6: CLI (`cmd/agent-cli/`)
- Phase 7~: HTTP API (`POST /v1/agent/run`)

---

## 다루지 않는 것 (Out of Scope)

### Agent 유형
- **브라우저 자동조작 Agent** — Playwright, Puppeteer 기반 웹 자동화
- **자율 코딩 Agent** — 코드 생성, 수정, 실행, 테스트까지 자동화하는 Agent
- **자율 배포 Agent** — CI/CD 파이프라인을 직접 트리거하는 Agent
- **특정 프레임워크 Wrapper** — LangChain, LangGraph 위에 얹는 구조

### 기능
- 멀티모달 입력 처리 (이미지, 오디오, 파일 업로드)
- 실시간 스트리밍 응답 (SSE, WebSocket)
- 사용자 인증 / 권한 관리
- 외부 SaaS 연동 (Slack, Notion, GitHub API 등 실제 연동)
- 프론트엔드 UI

### 운영
- Kubernetes 배포 (Phase 8 이전)
- 다중 인스턴스 수평 확장 (Phase 7까지는 단일 프로세스)
- 비용 청구 시스템

---

## 경계가 애매한 영역

| 항목 | 판단 기준 |
|------|----------|
| 실제 외부 API 연동 (날씨, 검색 등) | mock으로 대체. 실제 연동은 선택 사항 |
| pgvector / Qdrant | Phase 8 이후 선택적으로 도입. Phase 4는 태그 기반 검색으로 충분 |
| Kafka | Phase 7까지는 in-memory channel로 대체. Phase 8에서 선택 도입 |
| OpenTelemetry | Phase 8 대상. 그 전엔 structured log로 대체 |

---

## 이 문서를 갱신해야 하는 시점

- 새로운 Agent 유형을 추가하려 할 때
- 외부 시스템 연동 요구가 생겼을 때
- Phase 계획이 변경될 때
