.PHONY: build test test-unit test-integration lint vet

# Unit test only — integration 테스트는 제외 (//go:build integration 태그가 붙은 파일은 빌드되지 않음)
# CI(GitHub Actions)는 이 타겟만 실행한다.
test:
	go test ./...

# Unit test (test와 동일, CI 표준 타겟)
test-unit:
	go test ./...

# Integration test — Redis/Postgres가 실제로 떠 있어야 함 (docker-compose up 선행 필요)
test-integration:
	go test -tags integration ./...

# 전체 빌드 검증
build:
	go build ./...

# vet
vet:
	go vet ./...

# lint (golangci-lint 설치 필요)
lint:
	golangci-lint run ./...
