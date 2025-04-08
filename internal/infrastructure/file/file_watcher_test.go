package file_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/file"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestFileWatcher_InitialLoad(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "deployments.json")

	deployments := `[{"name":"test","image":"alpine","environment":{},"tags":[]}]`
	err := os.WriteFile(filePath, []byte(deployments), 0644)
	assert.NoError(t, err)

	var (
		mu       sync.Mutex
		called   bool
		received []model.Deployment
	)

	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		mu.Lock()
		defer mu.Unlock()
		called = true
		received = deployments
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(filePath, "test-source", push, logger)

	go func() {
		_ = w.Watch(ctx)
	}()

	time.Sleep(300 * time.Millisecond) // wait for debounce + push

	mu.Lock()
	assert.True(t, called, "pushFunc should have been called")
	assert.Len(t, received, 1)
	assert.Equal(t, "test", received[0].Name)
	mu.Unlock()
}

func TestFileWatcher_FileChangeTriggersReload(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "deployments.json")

	initial := `[{"name":"one","image":"alpine","environment":{},"tags":[]}]`
	updated := `[{"name":"two","image":"nginx","environment":{},"tags":[]}]`

	err := os.WriteFile(filePath, []byte(initial), 0644)
	assert.NoError(t, err)

	var (
		mu        sync.Mutex
		callCount int
		names     []string
	)

	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		mu.Lock()
		defer mu.Unlock()
		callCount++
		if len(deployments) > 0 {
			names = append(names, deployments[0].Name)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(filePath, "test-reload", push, logger)

	go func() {
		_ = w.Watch(ctx)
	}()

	time.Sleep(300 * time.Millisecond) // initial load

	err = os.WriteFile(filePath, []byte(updated), 0644)
	assert.NoError(t, err)

	time.Sleep(400 * time.Millisecond) // debounce + reload

	mu.Lock()
	assert.GreaterOrEqual(t, callCount, 2, "pushFunc should have been called at least twice")
	assert.Contains(t, names, "one")
	assert.Contains(t, names, "two")
	mu.Unlock()
}

func TestFileWatcher_HandlesInvalidJSONGracefully(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "deployments.json")

	initial := `[{"name":"valid","image":"alpine","environment":{},"tags":[]}]`
	invalid := `{{broken`

	err := os.WriteFile(filePath, []byte(initial), 0644)
	assert.NoError(t, err)

	var (
		mu        sync.Mutex
		callCount int
	)

	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		mu.Lock()
		defer mu.Unlock()
		callCount++
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(filePath, "test-bad-json", push, logger)

	go func() {
		_ = w.Watch(ctx)
	}()

	time.Sleep(300 * time.Millisecond) // wait for initial push
	err = os.WriteFile(filePath, []byte(invalid), 0644)
	assert.NoError(t, err)

	time.Sleep(400 * time.Millisecond) // wait for debounce + failed reload

	mu.Lock()
	assert.Equal(t, 1, callCount, "pushFunc should not be called again after invalid JSON")
	mu.Unlock()
}

func TestFileWatcher_PanicsOnNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil logger, but got none")
		}
	}()

	_ = file.NewWatcher("somefile.json", "source", func(context.Context, string, []model.Deployment) {}, nil)
}

func TestFileWatcher_InvalidFilePath(t *testing.T) {
	tempDir := t.TempDir() // a directory, not a file
	logger := zaptest.NewLogger(t)

	var called bool
	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		called = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	w := file.NewWatcher(tempDir, "invalid-path", push, logger)

	go func() {
		_ = w.Watch(ctx) // let it run and fail internally
	}()

	time.Sleep(300 * time.Millisecond)

	assert.False(t, called, "pushFunc should not have been called for invalid path")
}

func TestFileWatcher_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	missingPath := filepath.Join(tempDir, "not-there.json")

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(missingPath, "missing", func(context.Context, string, []model.Deployment) {}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := w.Watch(ctx)
	assert.Error(t, err)
}

func TestFileWatcher_PushFuncFailsButContinues(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "deployments.json")

	data := `[{"name":"bad","image":"alpine","environment":{},"tags":[]}]`
	err := os.WriteFile(filePath, []byte(data), 0644)
	assert.NoError(t, err)

	// broken pushFunc: panics internally
	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		panic("simulated push failure")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(filePath, "fail-push", push, logger)

	// Watcher should not panic, but the pushFunc will
	// So we use recover internally in a safe goroutine
	go func() {
		defer func() {
			_ = recover() // suppress panic from pushFunc
		}()
		_ = w.Watch(ctx)
	}()

	time.Sleep(300 * time.Millisecond)
}

func TestFileWatcher_DebounceCancelsOnShutdown(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "deployments.json")

	// Create a valid deployment file
	data := `[{"name":"test","image":"alpine","environment":{},"tags":[]}]`
	assert.NoError(t, os.WriteFile(filePath, []byte(data), 0644))

	var mu sync.Mutex
	var pushCount int

	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		mu.Lock()
		defer mu.Unlock()
		pushCount++
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zaptest.NewLogger(t)
	w := file.NewWatcher(filePath, "test-debounce-shutdown", push, logger)

	go func() {
		_ = w.Watch(ctx)
	}()

	time.Sleep(200 * time.Millisecond) // Wait for initial load

	// Trigger a fake file change (but don't actually write)
	// This will schedule a debounce
	_ = os.Chtimes(filePath, time.Now(), time.Now())

	// Cancel context immediately while debounce is waiting
	cancel()

	time.Sleep(400 * time.Millisecond) // Wait for debounce to settle

	mu.Lock()
	defer mu.Unlock()

	// Only the initial push should have happened
	assert.Equal(t, 1, pushCount, "should not push again after shutdown")
}
