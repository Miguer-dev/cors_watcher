package main

import (
	"os"
	"testing"

	"github.com/Miguer-dev/cors_watcher/internal/assert"
)

func createTempFile(t *testing.T, fileName string, fileContent string) *os.File {
	tempFile, err := os.CreateTemp("", fileName)
	assert.NilError(t, err)

	_, err = tempFile.WriteString(fileContent)
	assert.NilError(t, err)

	err = tempFile.Close()
	assert.NilError(t, err)

	return tempFile
}
