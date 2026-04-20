package memory

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/zipkero/agent-runtime/internal/types"
)

// memoryRepositorySuite 는 MemoryRepository 구현체에 공통으로 적용되는 테스트 케이스다.
// InMemoryMemoryRepository와 PostgresMemoryRepository 모두 동일한 케이스를 통과해야 한다.
func memoryRepositorySuite(t *testing.T, repo MemoryRepository) {
	t.Helper()
	ctx := context.Background()

	t.Run("Save_and_LoadByTags", func(t *testing.T) {
		m := types.Memory{
			ID:        "00000000-0000-0000-0000-000000000001",
			UserID:    "user-1",
			Content:   "test memory content",
			Tags:      []string{"go", "testing"},
			CreatedAt: time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC),
		}

		if err := repo.Save(ctx, m); err != nil {
			t.Fatalf("Save: %v", err)
		}

		got, err := repo.LoadByTags(ctx, []string{"go"}, 10)
		if err != nil {
			t.Fatalf("LoadByTags: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("LoadByTags: got %d records, want 1", len(got))
		}
		if got[0].ID != m.ID {
			t.Errorf("ID: got %q, want %q", got[0].ID, m.ID)
		}
		if got[0].Content != m.Content {
			t.Errorf("Content: got %q, want %q", got[0].Content, m.Content)
		}
	})

	t.Run("LoadByTags_OR_condition", func(t *testing.T) {
		// 서로 다른 태그를 가진 2개의 메모리를 저장한 뒤
		// 두 태그를 함께 조회하면 OR 조건으로 둘 다 반환되어야 한다.
		const m1ID = "00000000-0000-0000-0000-000000000011"
		const m2ID = "00000000-0000-0000-0000-000000000012"
		m1 := types.Memory{
			ID:        m1ID,
			UserID:    "user-1",
			Content:   "memory with alpha tag",
			Tags:      []string{"alpha"},
			CreatedAt: time.Date(2026, 4, 12, 1, 0, 0, 0, time.UTC),
		}
		m2 := types.Memory{
			ID:        m2ID,
			UserID:    "user-1",
			Content:   "memory with beta tag",
			Tags:      []string{"beta"},
			CreatedAt: time.Date(2026, 4, 12, 2, 0, 0, 0, time.UTC),
		}

		if err := repo.Save(ctx, m1); err != nil {
			t.Fatalf("Save m1: %v", err)
		}
		if err := repo.Save(ctx, m2); err != nil {
			t.Fatalf("Save m2: %v", err)
		}

		got, err := repo.LoadByTags(ctx, []string{"alpha", "beta"}, 10)
		if err != nil {
			t.Fatalf("LoadByTags: %v", err)
		}
		if len(got) < 2 {
			t.Fatalf("LoadByTags OR: got %d records, want at least 2", len(got))
		}

		ids := make(map[string]bool)
		for _, m := range got {
			ids[m.ID] = true
		}
		if !ids[m1ID] || !ids[m2ID] {
			t.Errorf("LoadByTags OR: expected both %s and %s, got IDs %v", m1ID, m2ID, ids)
		}
	})

	t.Run("LoadByTags_partial_tag_match", func(t *testing.T) {
		// 태그 중 하나만 일치해도 반환되어야 한다.
		const partialID = "00000000-0000-0000-0000-000000000021"
		m := types.Memory{
			ID:        partialID,
			UserID:    "user-1",
			Content:   "memory with multiple tags",
			Tags:      []string{"gamma", "delta"},
			CreatedAt: time.Date(2026, 4, 12, 3, 0, 0, 0, time.UTC),
		}

		if err := repo.Save(ctx, m); err != nil {
			t.Fatalf("Save: %v", err)
		}

		got, err := repo.LoadByTags(ctx, []string{"gamma", "nonexistent"}, 10)
		if err != nil {
			t.Fatalf("LoadByTags: %v", err)
		}

		found := false
		for _, g := range got {
			if g.ID == partialID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("LoadByTags partial: %s not found despite matching tag 'gamma'", partialID)
		}
	})

	t.Run("LoadByTags_empty_tags_returns_empty", func(t *testing.T) {
		got, err := repo.LoadByTags(ctx, []string{}, 10)
		if err != nil {
			t.Fatalf("LoadByTags empty tags: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("LoadByTags empty tags: got %d records, want 0", len(got))
		}
	})

	t.Run("LoadByTags_nil_tags_returns_empty", func(t *testing.T) {
		got, err := repo.LoadByTags(ctx, nil, 10)
		if err != nil {
			t.Fatalf("LoadByTags nil tags: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("LoadByTags nil tags: got %d records, want 0", len(got))
		}
	})

	t.Run("LoadByTags_limit_exceeded", func(t *testing.T) {
		// 동일 태그로 3개 저장 후 limit=2로 조회하면 2개만 반환되어야 한다.
		for i, id := range []string{
			"00000000-0000-0000-0000-000000000031",
			"00000000-0000-0000-0000-000000000032",
			"00000000-0000-0000-0000-000000000033",
		} {
			m := types.Memory{
				ID:        id,
				UserID:    "user-1",
				Content:   "limit test memory",
				Tags:      []string{"limitcheck"},
				CreatedAt: time.Date(2026, 4, 12, 10+i, 0, 0, 0, time.UTC),
			}
			if err := repo.Save(ctx, m); err != nil {
				t.Fatalf("Save %s: %v", id, err)
			}
		}

		got, err := repo.LoadByTags(ctx, []string{"limitcheck"}, 2)
		if err != nil {
			t.Fatalf("LoadByTags limit: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("LoadByTags limit: got %d records, want 2", len(got))
		}
	})

	t.Run("LoadByTags_no_match_returns_empty", func(t *testing.T) {
		got, err := repo.LoadByTags(ctx, []string{"nonexistent-tag-xyz"}, 10)
		if err != nil {
			t.Fatalf("LoadByTags no match: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("LoadByTags no match: got %d records, want 0", len(got))
		}
	})
}

func TestInMemoryMemoryRepository(t *testing.T) {
	repo := NewInMemoryMemoryRepository()
	memoryRepositorySuite(t, repo)
}

const testPostgresDSN = "postgres://agent:agent@localhost:5432/agent_runtime?sslmode=disable"

// newTestPostgresPool 은 localhost:5432 Postgres에 연결을 시도한다.
// Postgres를 사용할 수 없으면 테스트를 건너뛴다.
func newTestPostgresPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, testPostgresDSN)
	if err != nil {
		t.Skipf("Postgres unavailable (%v): skipping Postgres tests", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("Postgres ping failed (%v): skipping Postgres tests", err)
	}
	return pool
}

// cleanupMemoriesTable 은 테스트 시작 전 memories 테이블의 데이터를 정리한다.
func cleanupMemoriesTable(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := pool.Exec(ctx, "DELETE FROM memories"); err != nil {
		t.Fatalf("cleanup memories table: %v", err)
	}
}

func TestPostgresMemoryRepository(t *testing.T) {
	pool := newTestPostgresPool(t)
	t.Cleanup(func() { pool.Close() })

	// Migrate 실행 — database/sql 필요
	db, err := sql.Open("pgx", testPostgresDSN)
	if err != nil {
		t.Fatalf("open sql.DB for migrate: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	// 테스트 시작 전 데이터 정리
	cleanupMemoriesTable(t, pool)
	t.Cleanup(func() { cleanupMemoriesTable(t, pool) })

	repo := NewPostgresMemoryRepository(pool)
	memoryRepositorySuite(t, repo)
}

// TestPostgresMemoryRepository_OrderByCreatedAtDesc 는 LoadByTags 가 최신순으로 반환하는지 검증한다.
func TestPostgresMemoryRepository_OrderByCreatedAtDesc(t *testing.T) {
	pool := newTestPostgresPool(t)
	t.Cleanup(func() { pool.Close() })

	db, err := sql.Open("pgx", testPostgresDSN)
	if err != nil {
		t.Fatalf("open sql.DB for migrate: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	cleanupMemoriesTable(t, pool)
	t.Cleanup(func() { cleanupMemoriesTable(t, pool) })

	repo := NewPostgresMemoryRepository(pool)

	const oldID = "00000000-0000-0000-0000-000000000041"
	const newID = "00000000-0000-0000-0000-000000000042"
	// 오래된 것부터 저장
	old := types.Memory{
		ID:        oldID,
		UserID:    "user-1",
		Content:   "older memory",
		Tags:      []string{"order"},
		CreatedAt: time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
	}
	recent := types.Memory{
		ID:        newID,
		UserID:    "user-1",
		Content:   "newer memory",
		Tags:      []string{"order"},
		CreatedAt: time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC),
	}

	if err := repo.Save(ctx, old); err != nil {
		t.Fatalf("Save old: %v", err)
	}
	if err := repo.Save(ctx, recent); err != nil {
		t.Fatalf("Save recent: %v", err)
	}

	got, err := repo.LoadByTags(ctx, []string{"order"}, 10)
	if err != nil {
		t.Fatalf("LoadByTags: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("LoadByTags: got %d, want 2", len(got))
	}
	if got[0].ID != newID {
		t.Errorf("first result should be newest: got %q, want %q", got[0].ID, newID)
	}
	if got[1].ID != oldID {
		t.Errorf("second result should be oldest: got %q, want %q", got[1].ID, oldID)
	}
}
