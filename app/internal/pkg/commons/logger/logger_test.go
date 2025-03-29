package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLogger_NewLogger(t *testing.T) {
	logger := NewLogger()

	require.NotNil(t, logger)
}

func TestLogger_Attributes(t *testing.T) {
	var buf bytes.Buffer
	var mu sync.Mutex
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		io.Copy(&buf, r)
	}()

	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
		w.Close()
		wg.Wait()
	}()

	// Ensure log is fully written
	attributes := map[string]any{
		"user_id": "123",
		"email":   "test@example.com",
	}

	logger := NewLogger()
	logger.Info(context.Background()).Attributes("Created new user", attributes)

	// Small delay to ensure log is written
	time.Sleep(10 * time.Millisecond)

	var logMap map[string]any
	jsonData := buf.Bytes()
	err = json.Unmarshal(jsonData, &logMap)

	if err != nil {
		t.Logf("Failed to unmarshal JSON: %v", err)
		t.Logf("Captured JSON data: %s", string(jsonData))
		require.NoError(t, err, "Should parse log JSON without error")
	}

	require.Equal(t, "Created new user", logMap["msg"])
	require.Equal(t, "INFO", logMap["level"])
	require.Equal(t, "123", logMap["user_id"])
	require.Equal(t, "test@example.com", logMap["email"])
}
