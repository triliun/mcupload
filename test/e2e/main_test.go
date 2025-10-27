package e2e

import (
	"os"
	"testing"

	"github.com/triliun/mcupload/backend/logger"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	// Setup
	db, err := setup()
	if err != nil {
		logger.Fatal("Setup failed", zap.Error(err))
	}
	defer db.Close()

	// Run tests
	exitCode := m.Run()

	// Cleanup
	cleanup()

	os.Exit(exitCode)
}
