package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func BackupDatabase(dbPath, backupDir string) (string, error) {
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("pos_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupFile)

	return backupPath, copyFile(dbPath, backupPath)
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = destinationFile.Close()
	}()

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return destinationFile.Sync()
}
